package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
	"github.com/raduq/goboss/ops"
)

var (
	ps       *exec.Cmd
	deployed bool
	started  bool
)

type data struct {
	started  bool `json:"started"`
	deployed bool `json:"deployed"`
}

func main() {
	r := mux.NewRouter()
	r.Path("/").Methods(http.MethodGet).HandlerFunc(index)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.PathPrefix("/scripts/").Handler(http.StripPrefix("/scripts/", http.FileServer(http.Dir("scripts"))))
	r.HandleFunc("/goboss/status", status).Methods("GET")
	r.HandleFunc("/goboss/start", bossStart).Methods("POST")
	r.HandleFunc("/goboss/build", buildArtifact).Methods("POST")
	r.HandleFunc("/goboss/unbuild", unbuild).Methods("POST")
	r.HandleFunc("/goboss/kill", bossKill).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func index(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles("templates/index.html")).Execute(w, nil)
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, string(getData()))
}

func bossStart(w http.ResponseWriter, r *http.Request) {
	cmd, err := ops.Start(ops.NewConfig())
	if err != nil {
		fmt.Println(err)
	}
	ps = cmd
	started = true
}

func bossKill(w http.ResponseWriter, r *http.Request) {
	ps.Process.Kill()
	started = false
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
		deployed = true
	}
}

func unbuild(w http.ResponseWriter, r *http.Request) {
	config := ops.NewConfig()
	err := ops.RemoveContents(config.DeploymentFolder)
	if err != nil {
		fmt.Printf("Error on clear jboss deployment folder\n")
	} else {
		deployed = false
	}
}

func getData() []byte {
	var d data
	d.deployed = deployed
	d.started = started

	respose, err := json.Marshal(d)
	if err != nil {
		fmt.Println(err)
	}
	return respose
}
