package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/apex/log"
	"github.com/apex/log/handlers/logfmt"
	"github.com/raduq/goboss/config"
	"github.com/raduq/goboss/ops"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

var (
	ps  *exec.Cmd
	hot bool
	raw bool
)

func main() {
	cfg := config.MustGet()

	log.SetHandler(logfmt.Default)
	log.SetLevel(log.MustParseLevel(strings.ToLower(cfg.LogLevel)))
	log.Info("initializing")

	app := cli.NewApp()
	app.Name = "goboss"
	app.Usage = "Build and run Jboss Projects"
	app.Action = defaultAction(cfg)
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "hot, H",
			Usage:       "Only build the project and copy to deployments folder",
			Destination: &hot,
		},
		cli.BoolFlag{
			Name:        "raw, R",
			Usage:       "Only starts the server without building the project",
			Destination: &raw,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "clean",
			Aliases: []string{"c"},
			Usage:   "Delete all content of deployments folder",
			Action:  cleanAction(cfg.JbossHome),
		},
	}
	app.Version = "1.0.robot"
	err := app.Run(os.Args)
	if err != nil {
		log.WithError(err).Fatal("error starting goboss")
	}
}

func defaultAction(cfg config.Config) func(c *cli.Context) {
	return func(c *cli.Context) {
		buildArtifact(cfg.JbossHome, cfg.BuildArgs)
		bossStart(cfg.JbossHome, cfg.Args)
	}
}

func cleanAction(jbossHome string) func(c *cli.Context) {
	return func(c *cli.Context) {
		log.Info("cleaning deployments folder...")
		unbuild(jbossHome)
		log.Info("clear!")
	}
}

func bossStart(jbossHome string, args string) {
	if !hot {
		cmd, err := ops.Start(jbossHome, args)
		if err != nil {
			log.WithError(err).Fatal("error starting server")
		}
		ps = cmd
	}
}

func buildArtifact(jbossHome string, buildArgs []string) {
	if !raw {
		var err error

		unbuild(jbossHome)
		verifyMaven()

		var p struct {
			Dir      string `yaml:"dir"`
			Target   string `yaml:"target-dir"`
			Artifact string `yaml:"artifact"`
		}

		yamlFile, err := ioutil.ReadFile("projects.yml")
		if err != nil {
			log.Fatalf("cannot read projects.yml %q\n", err)
		}

		err = yaml.Unmarshal(yamlFile, &p)
		if err != nil {
			log.Fatalf("cannot unmarshal projects.yml %q\n", err)
		}

		c := "/usr/bin/mvn"
		log.Info("building " + p.Dir)
		ops.ExecuteAndPrint(p.Dir, c, buildArgs)

		log.Info("copying " + p.Target + "/" + p.Artifact + " to " + jbossHome + "/standalone/deployments/" + p.Artifact)
		err = ops.CopyFile(p.Target+"/"+p.Artifact, jbossHome+"/standalone/deployments/"+p.Artifact)
		if err != nil {
			log.Fatalf("file copy failed %q\n", err)
		} else {
			log.Info("file copy succeeded\n")
		}
	}
}

func verifyMaven() {
	path, err := exec.LookPath("mvn")
	if err != nil {
		log.Fatal("maven not found on PATH")
	}
	log.Infof("maven is available %s\n", path)
}

func unbuild(jbossHome string) {
	err := ops.RemoveContents(jbossHome + "/standalone/deployments/")
	if err != nil {
		log.Fatal("error on cleaning jboss deployments folder")
	}
}
