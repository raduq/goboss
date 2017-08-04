package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/raduq/goboss/ops"
)

func main() {
	art := "standalone/deployments/contaazul-app.ear"
	pom := "/home/raduansantos/git/ContaAzul"
	ca := "/home/raduansantos/git/ContaAzul/contaazul-app/target/contaazul-app.ear"
	jb := "/home/raduansantos/Dev/Server/jboss-contaazul-2.0.0-SNAPSHOT/"

	err := ops.RemoveContents(jb + "standalone/deployments")
	if err != nil {
		fmt.Printf("Error on clear jboss deployment folder\n")
	} else {
		buildCa(ca, jb, art, pom)
	}
}

func buildCa(ca, jb, art, pom string) {
	path, err := exec.LookPath("mvn")
	if err != nil {
		log.Fatal("Maven not found on PATH")
	}
	fmt.Printf("Maven is available %s\n", path)

	arguments := []string{"-T 1C", "package", "-o", "-Dmaven.test.skip", "-Dcheckstyle.skip", "-Denforcer.skip", "-Djacoco.skip"}
	ops.Execute(pom, "/usr/bin/mvn", arguments...)

	err = ops.CopyFile(ca, jb+art)
	if err != nil {
		fmt.Printf("CopyFile failed %q\n", err)
	} else {
		fmt.Printf("CopyFile succeeded\n")
		err = ops.Start(jb)
		if err != nil {
			fmt.Println(err)
		}
	}
}
