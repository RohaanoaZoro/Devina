package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/bitly/go-simplejson"
	guuid "github.com/google/uuid"
	"github.com/gorilla/mux"
)

var owner string = "lakbhat"
var repo string = "LanQuill"
var sha string = "dev_ui_editor"
var bearer = "Bearer " + ""
var url string = "https://api.github.com/repos/" + owner + "/" + repo + "/"
var since string = ""

var GitInitialPath string = "./Git"
var Credspath string = "./Credentials"

func sendRes(w http.ResponseWriter, jsonData *simplejson.Json) {

	// JSON encode jsonData
	payload, err := jsonData.MarshalJSON()
	if err != nil {
		log.Println(err, "\tstatus_code: 992")
		http.Error(w, "Internal Error", http.StatusMethodNotAllowed)
		return
	}

	// Return response JSON
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func ExecuteScriptsInDevina(shCmds string) (string, error) {

	var finalOutput string = ""
	cmd := exec.Command("/bin/sh", "-c", shCmds)
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
		return "Cannot Execute Cmd", nil
	}
	finalOutput += string(out)

	return finalOutput, nil
}

func DirExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func genUUID() string {
	id := guuid.New()
	fmt.Printf("github.com/google/uuid:         %s\n", id.String())

	return id.String()
}

func main() {

	log.Println("Devina Started")
	router := NewRouter()
	if err := http.ListenAndServe(":2007", router); err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}

	// tempScript := "cd " + "./GitDir/test" + "; ls; " + "git init; git remote add origin git@github.com:lakbhat/LanQuill.git;	git checkout --orphan dev_ui_editor;	git pull origin dev_ui_editor"

	// var finalOutput string = ""
	// cmd := exec.Command("/bin/sh", "-c", tempScript)
	// // cmd.Path = "./Git/test/"
	// out, err := cmd.Output()
	// if err != nil {
	// 	log.Println(err)
	// }
	// finalOutput += string(out)
	// log.Println("output", finalOutput)

	// var script string = "cd ./Git/PreProd/test; git status; git log; mkdir sup;"
	// ExecuteScriptsInDevina(script)
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	//Connection Related
	router.HandleFunc("/getconnparams", GetConnectionParams)
	router.HandleFunc("/postconnparams", PostConnParams)
	router.HandleFunc("/posteditconnparam", PostEditConnParamURL)

	//Branch Related
	router.HandleFunc("/getbranches", GetBranches)
	router.HandleFunc("/postcreatebranch", PostCreateBranch)
	router.HandleFunc("/getgithubbranches", GetGitHubBranches)

	//Integration Related
	router.HandleFunc("/getintegrations", GetIntegrations)
	router.HandleFunc("/postaddintegration", PostAddIntegration)

	//Commits Related
	router.HandleFunc("/getcommits", GetGithubCommitAPI)
	router.HandleFunc("/postconfirmcommits", PostConfirmCommits)

	return router
}
