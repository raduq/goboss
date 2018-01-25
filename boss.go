package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/raduq/goboss/ops"
	"github.com/urfave/cli"
)

var (
	ps  *exec.Cmd
	hot bool
	raw bool
)

func main() {
	app := cli.NewApp()
	app.Name = "goboss"
	app.Usage = "Build and run Jboss Projects"
	app.Action = defaultAction
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
			Action:  cleanAction,
		},
	}
	app.Version = "1.0.robot"
	app.Run(os.Args)
}

func defaultAction(c *cli.Context) {
	buildArtifact()
	bossStart()
}

func cleanAction(c *cli.Context) {
	fmt.Println("Cleaning deployments folder...")
	unbuild()
	fmt.Println("Done")
}

func bossStart() {
	if !hot {
		cmd, err := ops.Start(ops.NewConfig())
		if err != nil {
			fmt.Println(err)
		}
		ps = cmd
	}
}

func buildArtifact() {
	if !raw {
		var err error

		unbuild()
		verifyMaven()

		config := ops.NewConfig()
		ops.ExecuteAndPrint(config.ProjectFolder, config.Command, config.Arguments)

		pArtifact := config.TargetFolder + "/" + config.ArtifactName
		dArtifact := config.DeploymentFolder + "/" + config.ArtifactName
		err = ops.CopyFile(pArtifact, dArtifact)
		if err != nil {
			fmt.Printf("CopyFile failed %q\n", err)
		} else {
			fmt.Printf("CopyFile succeeded\n")
		}
	}
}

func verifyMaven() {
	path, err := exec.LookPath("mvn")
	if err != nil {
		log.Fatal("Maven not found on PATH")
	}
	fmt.Printf("Maven is available %s\n", path)
}

func unbuild() {
	config := ops.NewConfig()
	err := ops.RemoveContents(config.DeploymentFolder)
	if err != nil {
		fmt.Printf("Error on clear jboss deployment folder\n")
	}
}
