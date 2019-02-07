package envman

import (
	"errors"

	"github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/goinp/goinp"
	"gopkg.in/yaml.v2"
)

var (
	// CurrentEnvStoreFilePath ...
	CurrentEnvStoreFilePath string

	// ToolMode ...
	ToolMode bool
)

// -------------------
// --- Environment handling methods

// UpdateOrAddToEnvlist ...
func UpdateOrAddToEnvlist(oldEnvSlice []models.EnvironmentItemModel, newEnv models.EnvironmentItemModel, replace bool) ([]models.EnvironmentItemModel, error) {
	newKey, _, err := newEnv.GetKeyValuePair()
	if err != nil {
		return []models.EnvironmentItemModel{}, err
	}

	var newEnvs []models.EnvironmentItemModel
	exist := false

	if replace {
		match := 0
		for _, env := range oldEnvSlice {
			key, _, err := env.GetKeyValuePair()
			if err != nil {
				return []models.EnvironmentItemModel{}, err
			}

			if key == newKey {
				match = match + 1
			}
		}
		if match > 1 {
			if ToolMode {
				return []models.EnvironmentItemModel{}, errors.New("More then one env exist with key '" + newKey + "'")
			}
			msg := "   More then one env exist with key '" + newKey + "' replace all/append ['replace/append'] ?"
			answer, err := goinp.AskForString(msg)
			if err != nil {
				return []models.EnvironmentItemModel{}, err
			}

			switch answer {
			case "replace":
				break
			case "append":
				replace = false
				break
			default:
				return []models.EnvironmentItemModel{}, errors.New("Failed to parse answer: '" + answer + "' use ['replace/append']!")
			}
		}
	}

	for _, env := range oldEnvSlice {
		key, _, err := env.GetKeyValuePair()
		if err != nil {
			return []models.EnvironmentItemModel{}, err
		}

		if replace && key == newKey {
			exist = true
			newEnvs = append(newEnvs, newEnv)
		} else {
			newEnvs = append(newEnvs, env)
		}
	}

	if !exist {
		newEnvs = append(newEnvs, newEnv)
	}

	return newEnvs, nil
}

func removeDefaults(env *models.EnvironmentItemModel) error {
	opts, err := env.GetOptions()
	if err != nil {
		return err
	}

	if opts.Title != nil && *opts.Title == "" {
		opts.Title = nil
	}
	if opts.Description != nil && *opts.Description == "" {
		opts.Description = nil
	}
	if opts.Summary != nil && *opts.Summary == "" {
		opts.Summary = nil
	}
	if opts.IsRequired != nil && *opts.IsRequired == models.DefaultIsRequired {
		opts.IsRequired = nil
	}
	if opts.IsDontChangeValue != nil && *opts.IsDontChangeValue == models.DefaultIsDontChangeValue {
		opts.IsDontChangeValue = nil
	}
	if opts.IsTemplate != nil && *opts.IsTemplate == models.DefaultIsTemplate {
		opts.IsTemplate = nil
	}
	if opts.IsExpand != nil && *opts.IsExpand == models.DefaultIsExpand {
		opts.IsExpand = nil
	}
	if opts.IsSensitive != nil && *opts.IsSensitive == models.DefaultIsSensitive {
		opts.IsSensitive = nil
	}
	if opts.SkipIfEmpty != nil && *opts.SkipIfEmpty == models.DefaultSkipIfEmpty {
		opts.SkipIfEmpty = nil
	}

	(*env)[models.OptionsKey] = opts
	return nil
}

func generateFormattedYMLForEnvModels(envs []models.EnvironmentItemModel, unsets []string) (models.EnvsSerializeModel, error) {
	envMapSlice := []models.EnvironmentItemModel{}
	for _, env := range envs {
		err := removeDefaults(&env)
		if err != nil {
			return models.EnvsSerializeModel{}, err
		}

		hasOptions := false
		opts, err := env.GetOptions()
		if err != nil {
			return models.EnvsSerializeModel{}, err
		}

		if opts.Title != nil {
			hasOptions = true
		}
		if opts.Description != nil {
			hasOptions = true
		}
		if opts.Summary != nil {
			hasOptions = true
		}
		if len(opts.ValueOptions) > 0 {
			hasOptions = true
		}
		if opts.IsRequired != nil {
			hasOptions = true
		}
		if opts.IsDontChangeValue != nil {
			hasOptions = true
		}
		if opts.IsTemplate != nil {
			hasOptions = true
		}
		if opts.IsExpand != nil {
			hasOptions = true
		}
		if opts.IsSensitive != nil {
			hasOptions = true
		}
		if opts.SkipIfEmpty != nil {
			hasOptions = true
		}

		if !hasOptions {
			delete(env, models.OptionsKey)
		}

		envMapSlice = append(envMapSlice, env)
	}

	return models.EnvsSerializeModel{
		Envs:   envMapSlice,
		Unsets: unsets,
	}, nil
}

// -------------------
// --- File methods

// WriteEnvMapToFile ...
func WriteEnvMapToFile(pth string, envs []models.EnvironmentItemModel, unsets []string) error {
	if pth == "" {
		return errors.New("No path provided")
	}

	envYML, err := generateFormattedYMLForEnvModels(envs, unsets)
	if err != nil {
		return err
	}
	bytes, err := yaml.Marshal(envYML)
	if err != nil {
		return err
	}
	return fileutil.WriteBytesToFile(pth, bytes)
}

// InitAtPath ...
func InitAtPath(pth string) error {
	if exist, err := pathutil.IsPathExists(pth); err != nil {
		return err
	} else if !exist {
		if err := WriteEnvMapToFile(pth, []models.EnvironmentItemModel{}, []string{}); err != nil {
			return err
		}
	} else {
		errorMsg := "Path already exist: " + pth
		return errors.New(errorMsg)
	}
	return nil
}

// ParseEnvsYML ...
func ParseEnvsYML(bytes []byte) (models.EnvsSerializeModel, error) {
	var envstore models.EnvsSerializeModel
	if err := yaml.Unmarshal(bytes, &envstore); err != nil {
		return models.EnvsSerializeModel{}, err
	}
	for _, env := range envstore.Envs {
		if err := env.NormalizeValidateFillDefaults(); err != nil {
			return models.EnvsSerializeModel{}, err
		}
	}
	return envstore, nil
}

// ReadEnvs ...
func ReadEnvs(pth string) (models.EnvsSerializeModel, error) {
	bytes, err := fileutil.ReadBytesFromFile(pth)
	if err != nil {
		return models.EnvsSerializeModel{}, err
	}

	return ParseEnvsYML(bytes)
}

// ReadEnvsOrCreateEmptyList ...
func ReadEnvsOrCreateEmptyList() (models.EnvsSerializeModel, error) {
	envstore, err := ReadEnvs(CurrentEnvStoreFilePath)
	if err != nil {
		if err.Error() == "No environment variable list found" {
			err = InitAtPath(CurrentEnvStoreFilePath)
			return models.EnvsSerializeModel{}, err
		}
		return models.EnvsSerializeModel{}, err
	}
	return envstore, nil
}
