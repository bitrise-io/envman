package integration

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

// Golden parity suite: every case runs the binary under test (INTEGRATION_TEST_BINARY_PATH) in a hermetic
// work dir and HOME, and records exit code, stdout, stderr and every file change per step into testdata/golden/<case>.yml.
// Regenerate with: go test ./_tests/integration/ -run TestGolden -update
var updateGolden = flag.Bool("update", false, "update golden files")

const (
	goldenDir        = "testdata/golden"
	workPlaceholder  = "$WORK"
	maxInlineContent = 1024
)

type goldenStep struct {
	args  []string
	env   []string
	stdin *string
}

type goldenCase struct {
	name  string
	files map[string]string
	env   []string
	steps []goldenStep
}

type goldenFile struct {
	Steps []goldenResult `yaml:"steps"`
}

type goldenResult struct {
	Args         []string          `yaml:"args,flow"`
	Env          []string          `yaml:"env,omitempty"`
	Stdin        *string           `yaml:"stdin,omitempty"`
	ExitCode     int               `yaml:"exit_code"`
	Stdout       string            `yaml:"stdout,omitempty"`
	Stderr       string            `yaml:"stderr,omitempty"`
	ChangedFiles map[string]string `yaml:"changed_files,omitempty"`
	RemovedFiles []string          `yaml:"removed_files,omitempty"`
}

func envmanCmd(args ...string) goldenStep {
	return goldenStep{args: args}
}

func (s goldenStep) withStdin(stdin string) goldenStep {
	s.stdin = &stdin
	return s
}

func (s goldenStep) withEnv(env ...string) goldenStep {
	s.env = append(s.env, env...)
	return s
}

func TestGolden(t *testing.T) {
	bin, err := exec.LookPath(binPath())
	require.NoError(t, err)
	bin, err = filepath.Abs(bin)
	require.NoError(t, err)

	version, err := exec.Command(bin, "--version").Output()
	require.NoError(t, err)
	t.Logf("binary under test: %s (%s)", bin, strings.TrimSpace(string(version)))

	names := map[string]bool{}
	for _, tc := range goldenCases {
		require.False(t, names[tc.name], "duplicate case name: %s", tc.name)
		names[tc.name] = true

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := runGoldenCase(t, bin, strings.TrimSpace(string(version)), tc)
			pth := filepath.Join(goldenDir, tc.name+".yml")
			if *updateGolden {
				content, err := yaml.Marshal(goldenFile{Steps: got})
				require.NoError(t, err)
				require.NoError(t, os.MkdirAll(goldenDir, 0755))
				require.NoError(t, os.WriteFile(pth, content, 0644))
				return
			}

			content, err := os.ReadFile(pth)
			require.NoError(t, err, "missing golden file, run with -update")
			var want goldenFile
			require.NoError(t, yaml.Unmarshal(content, &want), pth)
			compareSteps(t, pth, want.Steps, got)
		})
	}
}

func compareSteps(t *testing.T, goldenPth string, want, got []goldenResult) {
	for i := 0; i < len(want) || i < len(got); i++ {
		var w, g goldenResult
		if i < len(want) {
			w = want[i]
		}
		if i < len(got) {
			g = got[i]
		}
		wantYML, gotYML := mustMarshal(t, w), mustMarshal(t, g)
		if wantYML == gotYML {
			continue
		}

		args := g.Args
		if i >= len(got) {
			args = w.Args
		}
		diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
			A:        difflib.SplitLines(wantYML),
			B:        difflib.SplitLines(gotYML),
			FromFile: "want (" + goldenPth + ")",
			ToFile:   "got",
			Context:  3,
		})
		t.Errorf("step %d/%d differs: envman %s\n%s", i+1, len(want), strings.Join(args, " "), diff)
	}
}

func mustMarshal(t *testing.T, v interface{}) string {
	b, err := yaml.Marshal(v)
	require.NoError(t, err)
	return string(b)
}

func runGoldenCase(t *testing.T, bin, version string, tc goldenCase) []goldenResult {
	work := t.TempDir()
	realWork, err := filepath.EvalSymlinks(work)
	require.NoError(t, err)
	home := filepath.Join(work, "home")
	require.NoError(t, os.MkdirAll(home, 0755))

	for name, content := range tc.files {
		pth := filepath.Join(work, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(pth), 0755))
		require.NoError(t, os.WriteFile(pth, []byte(content), 0644))
	}

	n := normalizer{work: work, realWork: realWork, version: version}
	resolve := func(s string) string { return strings.ReplaceAll(s, workPlaceholder, work) }

	var results []goldenResult
	prev := snapshotFiles(t, work)
	for _, step := range tc.steps {
		env := append(append([]string{}, tc.env...), step.env...)
		args := make([]string, len(step.args))
		for i, a := range step.args {
			args[i] = resolve(a)
		}

		cmd := exec.Command(bin, args...)
		cmd.Args[0] = "envman"
		cmd.Dir = work
		cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=" + home}
		for _, e := range env {
			cmd.Env = append(cmd.Env, resolve(e))
		}
		if step.stdin != nil {
			cmd.Stdin = strings.NewReader(*step.stdin)
		}
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		exitCode := 0
		if err := cmd.Run(); err != nil {
			var exitErr *exec.ExitError
			require.True(t, errors.As(err, &exitErr), "failed to run command: %s", err)
			exitCode = exitErr.ExitCode()
		}

		stdoutStr := n.normalize(trimTrailingSpaces(stdout.String()))
		if isPrintCommand(step.args) {
			stdoutStr = sortLines(stdoutStr)
		}

		result := goldenResult{
			Args:     append([]string{}, step.args...),
			Env:      env,
			Stdin:    step.stdin,
			ExitCode: exitCode,
			Stdout:   stdoutStr,
			Stderr:   n.normalize(trimTrailingSpaces(stderr.String())),
		}

		curr := snapshotFiles(t, work)
		for _, name := range sortedKeys(curr) {
			if prevContent, ok := prev[name]; !ok || prevContent != curr[name] {
				if result.ChangedFiles == nil {
					result.ChangedFiles = map[string]string{}
				}
				result.ChangedFiles[name] = fileSummary(curr[name])
			}
		}
		for _, name := range sortedKeys(prev) {
			if _, ok := curr[name]; !ok {
				result.RemovedFiles = append(result.RemovedFiles, name)
			}
		}
		prev = curr
		results = append(results, result)
	}

	return results
}

type normalizer struct {
	work, realWork, version string
}

var (
	ansiPattern         = regexp.MustCompile(`\x1b\[[0-9;]*m`)
	logrusTimePattern   = regexp.MustCompile(`(?m)^([A-Z]{4})\[\d{2}:\d{2}:\d{2}\]`)
	stdlogPrefixPattern = regexp.MustCompile(`(?m)^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} `)
	logrusTextTime      = regexp.MustCompile(`time="[^"]*"`)
	stackTracePattern   = regexp.MustCompile(`(?s)\ngoroutine \d+ \[.*`)
	panicPCPattern      = regexp.MustCompile(`pc=0x[0-9a-f]+`)
)

func (n normalizer) normalize(s string) string {
	s = stackTracePattern.ReplaceAllString(s, "\n<stack trace>\n")
	s = panicPCPattern.ReplaceAllString(s, "pc=<pc>")
	s = logrusTextTime.ReplaceAllString(s, `time="<time>"`)
	s = ansiPattern.ReplaceAllString(s, "")
	s = logrusTimePattern.ReplaceAllString(s, "$1[HH:MM:SS]")
	s = stdlogPrefixPattern.ReplaceAllString(s, "YYYY/MM/DD HH:MM:SS ")
	s = strings.ReplaceAll(s, n.realWork, workPlaceholder)
	s = strings.ReplaceAll(s, n.work, workPlaceholder)
	if n.version != "" {
		s = strings.ReplaceAll(s, n.version, "$VERSION")
	}
	return s
}

func isPrintCommand(args []string) bool {
	for i := 0; i < len(args); i++ {
		switch a := args[i]; {
		case a == "-l" || a == "--loglevel" || a == "-p" || a == "--path":
			i++
		case strings.HasPrefix(a, "-"):
		default:
			return a == "print" || a == "p"
		}
	}
	return false
}

// print's raw and envlist formats iterate a map, so their line order is random.
func sortLines(s string) string {
	if s == "" {
		return s
	}
	lines := strings.Split(strings.TrimSuffix(s, "\n"), "\n")
	sort.Strings(lines)
	return strings.Join(lines, "\n") + "\n"
}

func trimTrailingSpaces(s string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " ")
	}
	return strings.Join(lines, "\n")
}

func snapshotFiles(t *testing.T, root string) map[string]string {
	files := map[string]string{}
	err := filepath.WalkDir(root, func(pth string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		content, err := os.ReadFile(pth)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, pth)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(rel)] = string(content)
		return nil
	})
	require.NoError(t, err)
	return files
}

func fileSummary(content string) string {
	if len(content) <= maxInlineContent {
		return content
	}
	return fmt.Sprintf("<%d bytes, sha256 %x>", len(content), sha256.Sum256([]byte(content)))
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
