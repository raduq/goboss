package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/creamdog/gonfig"
	"github.com/raduq/goboss/ops"
)

type Config struct {
	SkipBuild bool
	Debug     bool
	Hot       bool
}

type ProjectDir struct {
	ProjectFolder string
	TargetFolder  string
	ArtifactName  string
}

type JbossDir struct {
	JbossFolder      string
	DeploymentFolder string
	LogsFolder       string
	BinFolder        string

	LogFile   string
	RunFile   string
	DebugFile string

	RunArgs string
}

type Build struct {
	Command   string
	Arguments []string
}

func main() {
	config := Config{false, false, false}
	verifyArgs(&config)

	f, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Cannot find config.json")
	}
	defer f.Close()

	gf, _ := gonfig.FromJson(f)

	pDir := ProjectDir{
		getProp(gf, "project.dir"),
		getProp(gf, "project.target.dir"),
		getProp(gf, "project.artifact")}

	jDir := JbossDir{
		getProp(gf, "jboss.dir"),
		getProp(gf, "jboss.deploy.dir"),
		getProp(gf, "jboss.log.dir"),
		getProp(gf, "jboss.bin.dir"),

		getProp(gf, "jboss.log.file"),
		getProp(gf, "jboss.run.file"),
		getProp(gf, "jboss.debug.file"),
		getProp(gf, "jboss.run.args")}

	build := Build{
		"/usr/bin/mvn",
		[]string{"-T 1C", "package", "-o", "-Dmaven.test.skip", "-Dcheckstyle.skip", "-Denforcer.skip", "-Djacoco.skip"}}

	err = ops.RemoveContents(jDir.DeploymentFolder)
	if err != nil {
		fmt.Printf("Error on clear jboss deployment folder\n")
	} else {
		buildArtifact(pDir, build, jDir, config)
	}
}

func buildArtifact(pDir ProjectDir, build Build, jDir JbossDir, config Config) {
	var err error
	if !config.SkipBuild {
		path, err := exec.LookPath("mvn")
		if err != nil {
			log.Fatal("Maven not found on PATH")
		}
		fmt.Printf("Maven is available %s\n", path)

		ops.Execute(pDir.ProjectFolder, build.Command, build.Arguments)
	}

	pArtifact := pDir.TargetFolder + "/" + pDir.ArtifactName
	dArtifact := jDir.DeploymentFolder + "/" + pDir.ArtifactName
	err = ops.CopyFile(pArtifact, dArtifact)
	if err != nil {
		fmt.Printf("CopyFile failed %q\n", err)
	} else {
		fmt.Printf("CopyFile succeeded\n")
		if !config.Hot {
			err = ops.Start(jDir.BinFolder, jDir.RunFile, jDir.DebugFile, jDir.RunArgs, jDir.LogsFolder, jDir.LogFile, config.Debug)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func getProp(gc gonfig.Gonfig, prop string) string {
	value, err := gc.GetString(prop, nil)
	if err != nil {
		log.Fatalf("Cannot find property %s", prop)
	}
	return value
}

func verifyArgs(conf *Config) {
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
}
