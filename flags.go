package main

import "github.com/codegangsta/cli"

const (
	PATH_KEY string = "path"
	P_KEY    string = "p"

	KEY_KEY string = "key"
	K_KEY   string = "k"

	VALUE_KEY string = "value"
	V_KEY     string = "v"

	VALUE_FILE_KEY string = "valuefile"
	VF_KEY         string = "vf"

	EXPAND_KEY string = "expand"
	E_KEY      string = "e"
)

var (
	// App flags
	flPath = cli.StringFlag{
		Name:  PATH_KEY + ", " + P_KEY,
		Value: "",
		Usage: "envman's working path, this is file path, with format {SOME_DIR/envstore.yml}",
	}
	flags = []cli.Flag{
		flPath,
	}

	// Command flags
	flKey = cli.StringFlag{
		Name:  KEY_KEY + ", " + K_KEY,
		Value: "",
		Usage: "key of the environment variable",
	}
	flValue = cli.StringFlag{
		Name:  VALUE_KEY + ", " + V_KEY,
		Value: "",
		Usage: "value of the environment variable",
	}
	flValueFile = cli.StringFlag{
		Name:  VALUE_FILE_KEY + ", " + VF_KEY,
		Value: "",
		Usage: "path of the environment variable value",
	}
	flIsExpand = cli.StringFlag{
		Name:  EXPAND_KEY + ", " + E_KEY,
		Value: "",
		Usage: "defines if should replace environment variables",
	}
)
