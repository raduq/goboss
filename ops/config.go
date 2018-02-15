package ops

import (
	"log"
	"os"

	"github.com/creamdog/gonfig"
)

//Config in the config.json file
type Config struct {
	Debug            bool
	Hot              bool
	SkipBuild        bool
	ProjectFolder    string
	TargetFolder     string
	ArtifactName     string
	JbossFolder      string
	DeploymentFolder string
	LogsFolder       string
	BinFolder        string
	LogFile          string
	RunFile          string
	DebugFile        string
	RunArgs          string
	Command          string
	Arguments        []string
}

//NewConfig created by reading the config.json file
func NewConfig() Config {
	config := verifyArgs()

	f, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Cannot find config.json")
	}
	defer f.Close()

	gf, _ := gonfig.FromJson(f)

	config.ProjectFolder = getProp(gf, "project.dir")
	config.TargetFolder = getProp(gf, "project.target.dir")
	config.ArtifactName = getProp(gf, "project.artifact")
	config.JbossFolder = getProp(gf, "jboss.dir")
	config.DeploymentFolder = getProp(gf, "jboss.deploy.dir")
	config.LogsFolder = getProp(gf, "jboss.log.dir")
	config.BinFolder = getProp(gf, "jboss.bin.dir")
	config.LogFile = getProp(gf, "jboss.log.file")
	config.RunFile = getProp(gf, "jboss.run.file")
	config.DebugFile = getProp(gf, "jboss.debug.file")
	config.RunArgs = getProp(gf, "jboss.run.args")
	config.Command = "/usr/bin/mvn"
	config.Arguments = []string{"-T 1C", "package", "-o", "-Dmaven.test.skip", "-Dcheckstyle.skip", "-Denforcer.skip", "-Djacoco.skip"}
	config.Debug = true

	return config
}

func verifyArgs() Config {
	var conf Config
	for _, argument := range os.Args {
		if argument == "debug" {
			conf.Debug = true
		}
		if argument == "skipBuild" {
			conf.SkipBuild = true
		}
		if argument == "hot" {
			conf.Hot = true
		}
	}
	return conf
}

func getProp(gc gonfig.Gonfig, prop string) string {
	value, err := gc.GetString(prop, nil)
	if err != nil {
		log.Fatalf("Cannot find property %s", prop)
	}
	return value
}
