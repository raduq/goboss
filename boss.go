package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
	"github.com/raduq/goboss/ops"
)

var ps Command

func main() {
	r := mux.NewRouter()
	r.Path("/").
		Methods(http.MethodGet).
		HandlerFunc(index)
	r.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/goboss/start", bossStart).Methods("GET")
	r.HandleFunc("/goboss/build", buildArtifact).Methods("GET")
	r.HandleFunc("/goboss/unbuild", unbuild).Methods("GET")
	r.HandleFunc("/goboss/kill", bossKill).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func index(w http.ResponseWriter, r *http.Request) (ps *Cmd) {
	template.Must(template.ParseFiles("templates/index.html")).Execute(w, nil)
}

func bossStart(w http.ResponseWriter, r *http.Request) {
	ps, err := ops.Start(ops.NewConfig())
	if err != nil {
		fmt.Println(err)
	}
}

func bossKill(w http.ResponseWriter, r *http.Request) {
	exec.Command("kill -9 ", ps.Process.pid)
}

func buildArtifact(w http.ResponseWriter, r *http.Request) {
	var err error

	unbuild(w, r)

	path, err := exec.LookPath("mvn")
	if err != nil {
		log.Fatal("Maven not found on PATH")
	}
	fmt.Printf("Maven is available %s\n", path)

	config := ops.NewConfig()
	ops.Execute(config.ProjectFolder, config.Command, config.Arguments)

	pArtifact := config.TargetFolder + "/" + config.ArtifactName
	dArtifact := config.DeploymentFolder + "/" + config.ArtifactName
	err = ops.CopyFile(pArtifact, dArtifact)
	if err != nil {
		fmt.Printf("CopyFile failed %q\n", err)
	} else {
		fmt.Printf("CopyFile succeeded\n")
	}
}

func unbuild(w http.ResponseWriter, r *http.Request) {
	config := ops.NewConfig()
	err := ops.RemoveContents(config.DeploymentFolder)
	if err != nil {
		fmt.Printf("Error on clear jboss deployment folder\n")
	}
}
