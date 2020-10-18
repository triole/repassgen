package main

import (
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	BUILDTAGS      string
	appName        = "pwgen"
	appMainVersion = "0.1"
	appDescription = "regex scheme password password generator"

	app         = kingpin.New(appName, appDescription)
	argsPattern = app.Arg("pattern", "regex pattern to generate passwords").Default("[0-9a-zA-Z:punct:]{32}").String()
	argsEntropy = app.Flag("entropy", "show pattern entropy").Short('e').Default("false").Bool()
	argsNo      = app.Flag("no", "number of passwords to generate").Short('n').Default("32").Int()
)

func argparse() {
	env := tEnv{
		Name:        appName,
		MainVersion: appMainVersion,
		Description: appDescription,
	}
	app.Version(makeInfoString(env, parseBuildtags(BUILDTAGS)))
	app.VersionFlag.Short('V')
	app.HelpFlag.Short('h')

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
