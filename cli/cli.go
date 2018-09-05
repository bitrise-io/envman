package cli

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/envman/envman"
	"github.com/bitrise-io/envman/version"
	"github.com/urfave/cli"
)

const (
	defaultEnvStoreName = ".envstore.yml"
	helpTemplate        = `
	NAME: {{.Name}} - {{.Usage}}
	
	USAGE: {{.Name}} {{if .Flags}}[OPTIONS] {{end}}COMMAND [arg...]
	
	VERSION: {{.Version}}{{if or .Author .Email}}
	
	AUTHOR:{{if .Author}}
	  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
	  {{.Email}}{{end}}{{end}}
	{{if .Flags}}
	GLOBAL OPTIONS:
	  {{range .Flags}}{{.}}
	  {{end}}{{end}}
	COMMANDS:
	  {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
	  {{end}}
	COMMAND HELP: {{.Name}} COMMAND --help/-h
	
	`
)

func before(c *cli.Context) error {
	// Init logging
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})

	logLevel, err := log.ParseLevel(c.String(LogLevelKey))
	if err != nil {
		return fmt.Errorf("Failed to parse log level: %s", err)
	}
	log.SetLevel(logLevel)

	// Ensure the envstore path
	envman.CurrentEnvStoreFilePath = c.String(PathKey)
	if envman.CurrentEnvStoreFilePath == "" {
		path, err := filepath.Abs(path.Join("./", defaultEnvStoreName))
		if err != nil {
			log.Fatal("[ENVMAN] - Failed to set envman work path in current dir:", err)
		} else {
			envman.CurrentEnvStoreFilePath = path
		}
	}

	envman.ToolMode = c.Bool(ToolKey)
	if envman.ToolMode {
		log.Info("[ENVMAN] - Tool mode on")
	}

	if _, err := envman.GetConfigs(); err != nil {
		log.Fatal("[ENVMAN] - Failed to init configs:", err)
	}

	return nil
}

// Run the Envman CLI.
func Run() {
	cli.HelpFlag = cli.BoolFlag{Name: HelpKey + ", " + helpKeyShort, Usage: "Show help."}
	cli.AppHelpTemplate = helpTemplate

	cli.VersionFlag = cli.BoolFlag{Name: VersionKey + ", " + versionKeyShort, Usage: "Print the version."}
	cli.VersionPrinter = func(c *cli.Context) { fmt.Println(c.App.Version) }

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "Environment variable manager"
	app.Version = version.VERSION

	app.Author = ""
	app.Email = ""

	app.Before = before

	app.Flags = flags
	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Fatal("[ENVMAN] - Finished:", err)
	}
}
