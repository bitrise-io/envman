package integration

import (
	"github.com/bitrise-io/envman/env"
	"github.com/bitrise-io/envman/models"
)

var SharedTestCases = []struct {
	Name string
	Envs []models.EnvironmentItemModel
	Want []env.Command
}{
	{
		Name: "empty env list",
		Envs: []models.EnvironmentItemModel{},
		Want: []env.Command{},
	},
	{
		Name: "unset env",
		Envs: []models.EnvironmentItemModel{
			{"A": "B", "opts": map[string]interface{}{"unset": true}},
		},
		Want: []env.Command{
			{Action: env.UnsetAction, Variable: env.Variable{Key: "A"}},
		},
	},
	{
		Name: "set env",
		Envs: []models.EnvironmentItemModel{
			{"A": "B", "opts": map[string]interface{}{}},
		},
		Want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "A", Value: "B"}},
		},
	},
	{
		Name: "set multiple envs",
		Envs: []models.EnvironmentItemModel{
			{"A": "B", "opts": map[string]interface{}{}},
			{"B": "C", "opts": map[string]interface{}{}},
		},
		Want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "A", Value: "B"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "B", Value: "C"}},
		},
	},
	{
		Name: "set int env",
		Envs: []models.EnvironmentItemModel{
			{"A": 12, "opts": map[string]interface{}{}},
		},
		Want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "A", Value: "12"}},
		},
	},
	{
		Name: "skip env",
		Envs: []models.EnvironmentItemModel{
			{"A": "B", "opts": map[string]interface{}{}},
			{"S": "", "opts": map[string]interface{}{"skip_if_empty": true}},
		},
		Want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "A", Value: "B"}},
			{Action: env.SkipAction, Variable: env.Variable{Key: "S"}},
		},
	},
	{
		Name: "skip env, do not skip if not empty",
		Envs: []models.EnvironmentItemModel{
			{"A": "B", "opts": map[string]interface{}{}},
			{"S": "T", "opts": map[string]interface{}{"skip_if_empty": true}},
		},
		Want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "A", Value: "B"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "S", Value: "T"}},
		},
	},
	{
		Name: "Env does only depend on envs declared before them",
		Envs: []models.EnvironmentItemModel{
			{"simulator_device": "$simulator_major", "opts": map[string]interface{}{"is_expand": true}},
			{"simulator_major": "12", "opts": map[string]interface{}{"is_expand": false}},
			{"simulator_os_version": "$simulator_device", "opts": map[string]interface{}{"is_expand": true}},
		},
		Want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_device", Value: ""}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_major", Value: "12"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_os_version", Value: ""}},
		},
	},
	{
		Name: "Env does only depend on envs declared before them (input order switched)",
		Envs: []models.EnvironmentItemModel{
			{"simulator_device": "$simulator_major", "opts": map[string]interface{}{"is_expand": true}},
			{"simulator_os_version": "$simulator_device", "opts": map[string]interface{}{"is_sensitive": false}},
			{"simulator_major": "12", "opts": map[string]interface{}{"is_expand": true}},
		},
		Want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_device", Value: ""}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_os_version", Value: ""}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_major", Value: "12"}},
		},
	},
	{
		Name: "Env does only depend on envs declared before them, envs in a loop",
		Envs: []models.EnvironmentItemModel{
			{"A": "$C", "opts": map[string]interface{}{"is_expand": true}},
			{"B": "$A", "opts": map[string]interface{}{"is_expand": true}},
			{"C": "$B", "opts": map[string]interface{}{"is_expand": true}},
		},
		Want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "A", Value: ""}},
			{Action: env.SetAction, Variable: env.Variable{Key: "B", Value: ""}},
			{Action: env.SetAction, Variable: env.Variable{Key: "C", Value: ""}},
		},
	},
	{
		Name: "Do not expand env if is_expand is false",
		Envs: []models.EnvironmentItemModel{
			{"SIMULATOR_OS_VERSION": "13.3", "opts": map[string]interface{}{"is_expand": true}},
			{"simulator_os_version": "$SIMULATOR_OS_VERSION", "opts": map[string]interface{}{"is_expand": false}},
		},
		Want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "SIMULATOR_OS_VERSION", Value: "13.3"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_os_version", Value: "$SIMULATOR_OS_VERSION"}},
		},
	},
	{
		Name: "Expand env, self reference",
		Envs: []models.EnvironmentItemModel{
			{"SIMULATOR_OS_VERSION": "$SIMULATOR_OS_VERSION", "opts": map[string]interface{}{"is_expand": true}},
		},
		Want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "SIMULATOR_OS_VERSION", Value: ""}},
		},
	},
	{
		Name: "Expand env, input contains env var",
		Envs: []models.EnvironmentItemModel{
			{"SIMULATOR_OS_VERSION": "13.3", "opts": map[string]interface{}{"is_expand": false}},
			{"simulator_os_version": "$SIMULATOR_OS_VERSION", "opts": map[string]interface{}{"is_expand": true}},
		},
		Want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "SIMULATOR_OS_VERSION", Value: "13.3"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_os_version", Value: "13.3"}},
		},
	},
	{
		Name: "Multi level env var expansion",
		Envs: []models.EnvironmentItemModel{
			{"A": "1", "opts": map[string]interface{}{"is_expand": true}},
			{"B": "$A", "opts": map[string]interface{}{"is_expand": true}},
			{"C": "prefix $B", "opts": map[string]interface{}{"is_expand": true}},
		},
		Want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "A", Value: "1"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "B", Value: "1"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "C", Value: "prefix 1"}},
		},
	},
	{
		Name: "Multi level env var expansion 2",
		Envs: []models.EnvironmentItemModel{
			{"SIMULATOR_OS_MAJOR_VERSION": "13", "opts": map[string]interface{}{"is_expand": true}},
			{"SIMULATOR_OS_MINOR_VERSION": "3", "opts": map[string]interface{}{"is_expand": true}},
			{"SIMULATOR_OS_VERSION": "$SIMULATOR_OS_MAJOR_VERSION.$SIMULATOR_OS_MINOR_VERSION", "opts": map[string]interface{}{"is_expand": true}},
			{"simulator_os_version": "$SIMULATOR_OS_VERSION", "opts": map[string]interface{}{"is_expand": true}},
		},
		Want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "SIMULATOR_OS_MAJOR_VERSION", Value: "13"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "SIMULATOR_OS_MINOR_VERSION", Value: "3"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "SIMULATOR_OS_VERSION", Value: "13.3"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_os_version", Value: "13.3"}},
		},
	},
	{
		Name: "Env expand, duplicate env declarations",
		Envs: []models.EnvironmentItemModel{
			{"simulator_os_version": "12.1", "opts": map[string]interface{}{}},
			{"simulator_device": "iPhone 8 ($simulator_os_version)", "opts": map[string]interface{}{"is_expand": "true"}},
			{"simulator_os_version": "13.3", "opts": map[string]interface{}{}},
		},
		Want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_os_version", Value: "12.1"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_device", Value: "iPhone 8 (12.1)"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_os_version", Value: "13.3"}},
		},
	},
	{
		Name: "Secrets inputs are marked as sensitive",
		Envs: []models.EnvironmentItemModel{
			{"simulator_os_version": "13.3", "opts": map[string]interface{}{"is_sensitive": false}},
			{"secret_input": "top secret", "opts": map[string]interface{}{"is_sensitive": true}},
		},
		Want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_os_version", Value: "13.3"}},
			{Action: env.SetAction, Variable: env.Variable{Key: "secret_input", Value: "top secret", IsSensitive: true}},
		},
	},
	{
		Name: "Input referencing secret env is marked as sensitive",
		Envs: []models.EnvironmentItemModel{
			{"SECRET_ENV": "top secret", "opts": map[string]interface{}{"is_sensitive": true}},
			{"simulator_device": "iPhone $SECRET_ENV", "opts": map[string]interface{}{"is_expand": true, "is_sensitive": false}},
		},
		Want: []env.Command{
			{Action: env.SetAction, Variable: env.Variable{Key: "SECRET_ENV", Value: "top secret", IsSensitive: true}},
			{Action: env.SetAction, Variable: env.Variable{Key: "simulator_device", Value: "iPhone top secret", IsSensitive: true}},
		},
	},
}
