package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/raduq/goboss/ops"
)

type Config struct {
	SkipBuild bool
	Debug     bool
	Hot       bool
}

type Targets struct {
	ProjectFolder    string
	ArtifactName     string
	ArtifactFolder   string
	JbossFolder      string
	DeploymentFolder string
}

type Logs struct {
	LogsFolder string
	LogFile    string
}

type Build struct {
	Command   string
	Arguments []string
}

type Bin struct {
	binFolder string
	runFile   string
	debugFile string
	runArgs   string
}

func main() {
	config := Config{false, false, false}
	verifyArgs(&config)

	targets := Targets{
		"/home/raduansantos/git/ContaAzul",
		"contaazul-app.ear",
		"/home/raduansantos/git/ContaAzul/contaazul-app/target",
		"/home/raduansantos/Dev/Server/jboss-contaazul-2.0.0-SNAPSHOT",
		"/home/raduansantos/Dev/Server/jboss-contaazul-2.0.0-SNAPSHOT/standalone/deployments"}

	bin := Bin{
		"/home/raduansantos/Dev/Server/jboss-contaazul-2.0.0-SNAPSHOT/bin",
		"standalone.sh",
		"debug.sh",
		" -b localhost --server-config=standalone.xml -Djboss.server.base.dir=/home/raduansantos/Dev/Server/jboss-contaazul-2.0.0-SNAPSHOT/standalone -P=/home/raduansantos/Dev/Server/jboss-contaazul-2.0.0-SNAPSHOT/standalone/configuration/contaazul.properties"}

	logs := Logs{
		"/home/raduansantos/Dev/Server/jboss-contaazul-2.0.0-SNAPSHOT/standalone/log",
		"server.log"}

	build := Build{
		"/usr/bin/mvn",
		[]string{"-T 1C", "package", "-o", "-Dmaven.test.skip", "-Dcheckstyle.skip", "-Denforcer.skip", "-Djacoco.skip"}}

	// art := "/contaazul-app.ear"
	// pom := "/home/raduansantos/git/ContaAzul"
	// ca := "/home/raduansantos/git/ContaAzul/contaazul-app/target/contaazul-app.ear"
	// jb := "/home/raduansantos/Dev/Server/jboss-contaazul-2.0.0-SNAPSHOT/"

	err := ops.RemoveContents(targets.DeploymentFolder)
	if err != nil {
		fmt.Printf("Error on clear jboss deployment folder\n")
	} else {
		buildArtifact(targets, build, bin, logs, config)
	}
}

func buildArtifact(targets Targets, build Build, bin Bin, logs Logs, config Config) {
	var err error
	if !config.SkipBuild {
		path, err := exec.LookPath("mvn")
		if err != nil {
			log.Fatal("Maven not found on PATH")
		}
		fmt.Printf("Maven is available %s\n", path)

		ops.Execute(targets.ProjectFolder, build.Command, build.Arguments)
	}

	pomFolder := targets.ArtifactFolder + "/" + targets.ArtifactName
	destination := targets.DeploymentFolder + "/" + targets.ArtifactName
	err = ops.CopyFile(pomFolder, destination)
	if err != nil {
		fmt.Printf("CopyFile failed %q\n", err)
	} else {
		fmt.Printf("CopyFile succeeded\n")
		if !config.Hot {
			err = ops.Start(bin.binFolder, bin.runFile, bin.debugFile, bin.runArgs,
				logs.LogsFolder, logs.LogFile, config.Debug)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
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
