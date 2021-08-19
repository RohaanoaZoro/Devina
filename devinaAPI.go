package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bitly/go-simplejson"
)

type IntegrationStruct struct {
	IntegrationId   string `json:"integrationId"`
	Connid          string `json:"connid"`
	IntegrationName string `json:"integrationName"`
	IbranchId       string `json:"ibranchId"`
	FbranchId       string `json:"fbranchId"`
	Iscript         string `json:"iscript"`
	Fscript         string `json:"fscript"`
	SSHscript       string `json:"sshscript"`
}

type confirmCommitStruct struct {
	IntegrationId string `json:"integrationid"`
	ConnId        string `json:"connid"`
	Commits       string `json:"commits"`
}

type BranchStruct struct {
	BranchId     string `json:"branchid"`
	ConnId       string `json:"connid"`
	BranchName   string `json:"branchname"`
	BranchScript string `json:"branchscript"`
}

func GetBranches(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Create empty return JSON
	jsonData := simplejson.New()

	var connId string = r.FormValue("connId")

	request_code, request_msg, BranchArr := MySQL_GetBranches(connId)

	jsonData.Set("status", request_msg)
	jsonData.Set("status_code", request_code)
	if request_code != 0 {
		log.Println("BranchArr", BranchArr)
		if len(BranchArr) > 0 {
			jsonData.Set("data", BranchArr)
		} else {
			jsonData.Set("data", []string{})
		}
	}
	sendRes(w, jsonData)

	return
}

func GetIntegrations(w http.ResponseWriter, r *http.Request) {

	log.Println("GetIntegrations", r.Method)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Create empty return JSON
	jsonData := simplejson.New()

	request_code, request_msg, BranchArr := MySQL_GetIntegrations()

	jsonData.Set("status", request_msg)
	jsonData.Set("status_code", request_code)
	if BranchArr != nil {
		jsonData.Set("data", BranchArr)
	}
	sendRes(w, jsonData)

	return
}

func PostConfirmCommits(w http.ResponseWriter, r *http.Request) {

	// Accept only POST
	if r.Method != "POST" {
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create empty return JSON
	jsonData := simplejson.New()

	var postBody confirmCommitStruct
	// Decode post body
	err := json.NewDecoder(r.Body).Decode(&postBody)
	if err != nil {
		log.Println("Error in UnMarsal")
		jsonData.Set("status", "Error")
		jsonData.Set("status_code", "990")
		jsonData.Set("message", "Cannot Unmarshal Post Body")
		sendRes(w, jsonData)
		return
	}

	log.Println("PostBody", postBody.IntegrationId)

	req_code, req_msg, scripts := MySQL_GetScripts(postBody.IntegrationId)
	if req_code == 0 {
		log.Println("Error in Getting Scripts SQL")
		jsonData.Set("status", "Error")
		jsonData.Set("status_code", req_code)
		jsonData.Set("message", req_msg)
		sendRes(w, jsonData)
		return
	}

	req_code, req_msg, creds := MYSQL_GetSingleConnectionParams(postBody.ConnId)
	if req_code == 0 {
		log.Println("Error in Getting Connection Params SQL")
		jsonData.Set("status", "Error")
		jsonData.Set("status_code", req_code)
		jsonData.Set("message", req_msg)
		sendRes(w, jsonData)
		return
	}

	// log.Println("Scripts", scripts)
	log.Println("creds", creds)

	var devinaScript string = strings.Replace(scripts[0], "$gitmsg", "\""+postBody.Commits+"\"", -1)
	log.Println("devinaScript", devinaScript)

	output, err := ExecuteScriptsInDevina(devinaScript)
	if err != nil {
		log.Println("Error in Devina Execute", err)
	}

	// output, err = ExecuteScriptInSSHServer(creds.Address, creds.Host, creds.PrivateKeyPath, scripts[1])
	// if err != nil {
	// 	log.Println("Error in SSH Execute", err)
	// }

	log.Println("output", output)

	jsonData.Set("status", req_msg)
	jsonData.Set("status_code", req_code)
	jsonData.Set("scripts", scripts)
	sendRes(w, jsonData)
}

func PostAddIntegration(w http.ResponseWriter, r *http.Request) {

	// Accept only POST
	if r.Method != "POST" {
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create empty return JSON
	jsonData := simplejson.New()

	var postBody IntegrationStruct
	// Decode post body
	err := json.NewDecoder(r.Body).Decode(&postBody)
	if err != nil {
		log.Println("Error in UnMarsal")
		jsonData.Set("status", "Error")
		jsonData.Set("status_code", "990")
		jsonData.Set("message", "Cannot Unmarshal Post Body")
		sendRes(w, jsonData)
		return
	}

	var IntegrationId string = genUUID()
	res_code, res_msg := MYSQL_AddIntegration(IntegrationId, postBody.Connid, postBody.IntegrationName, postBody.IbranchId, postBody.FbranchId, postBody.Iscript, postBody.Fscript, postBody.SSHscript)
	if res_code == 0 {
		log.Println("Error in SQL Insert", res_msg)
		jsonData.Set("status", "Error")
		jsonData.Set("status_code", "990")
		jsonData.Set("message", "Error in Adding script to SQL")
		sendRes(w, jsonData)
		return
	}

	jsonData.Set("status", "Success")
	jsonData.Set("status_code", "01")
	sendRes(w, jsonData)
}

func GetGithubCommitAPI(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create empty return JSON
	jsonData := simplejson.New()

	var branch string = r.FormValue("branch")
	var since string = r.FormValue("since")
	var per_page string = r.FormValue("per_page")

	if branch == "" {
		jsonData.Set("status", "Error")
		jsonData.Set("status_code", "02")
		jsonData.Set("message", "Missing Branch Field")
		sendRes(w, jsonData)
		return
	}

	var Commits []RequiredCommitInfo
	var res_code int
	var res_msg string
	if since == "" {
		res_code, res_msg, Commits = getCommits(branch, per_page)
	} else {
		res_code, res_msg, Commits = getCommitsByTime(branch, per_page, since)
	}
	if res_code != 1 {
		jsonData.Set("status", "Error")
		jsonData.Set("status_code", res_code)
		jsonData.Set("message", res_msg)
		sendRes(w, jsonData)
		return
	}

	jsonData.Set("status", "Sucess")
	jsonData.Set("status_code", "01")
	jsonData.Set("data", Commits)
	sendRes(w, jsonData)
}

func GetGitHubBranches(w http.ResponseWriter, r *http.Request) {

	log.Println("GetGitHubBranches")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create empty return JSON
	jsonData := simplejson.New()

	// var branch string = r.FormValue("repoName")
	// if branch == "" {
	// 	jsonData.Set("status", "Error")
	// 	jsonData.Set("status_code", "02")
	// 	jsonData.Set("message", "Missing Branch Field")
	// 	sendRes(w, jsonData)
	// 	return
	// }

	res_code, res_msg, Branches := GetGitBranches()
	if res_code != 1 {
		jsonData.Set("status", "Error")
		jsonData.Set("status_code", res_code)
		jsonData.Set("message", res_msg)
		sendRes(w, jsonData)
		return
	}

	log.Println("Branches", Branches)

	jsonData.Set("status", "Sucess")
	jsonData.Set("status_code", "01")
	jsonData.Set("data", Branches)
	sendRes(w, jsonData)
}

func GetConnectionParams(w http.ResponseWriter, r *http.Request) {

	log.Println("GetConnectionParams")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create empty return JSON
	jsonData := simplejson.New()

	res_code, res_msg, connParam := MYSQL_GetAllConnections()
	if res_code == 0 {
		log.Println(res_msg)
		jsonData.Set("status", "Error")
		jsonData.Set("status_code", "990")
		jsonData.Set("message", res_msg)
		sendRes(w, jsonData)
		return
	}

	jsonData.Set("status", "Success")
	jsonData.Set("status_code", 1)
	jsonData.Set("data", connParam)
	sendRes(w, jsonData)
}

func PostConnParams(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	var host string = r.FormValue("host")
	var address string = r.FormValue("address")
	var connname string = r.FormValue("name")
	var ConnId string = genUUID()

	jsonData := simplejson.New()

	// Get handler for filename, size and headers
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File", err)
		jsonData.Set("status", "Error")
		jsonData.Set("status_code", 0)
		jsonData.Set("message", err.Error())
		sendRes(w, jsonData)
		return
	}

	defer file.Close()

	var connPath string = Credspath + "/" + ConnId

	// This is path which we want to store the file
	if !DirExists(connPath) {
		os.MkdirAll(connPath, 0777)
	}

	f, err := os.OpenFile(connPath+"/privatekey.pem", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Println("errror in Upload file", err)
		jsonData.Set("status", "Error")
		jsonData.Set("status_code", 0)
		jsonData.Set("message", err.Error())
		sendRes(w, jsonData)
		return
	}

	// Copy the file to the destination path
	io.Copy(f, file)

	err = AddConnParameters(host, address, connPath+"/privatekey.pem", connname, ConnId)
	if err != nil {
		log.Println("errror in SQL Add Params", err)
		jsonData.Set("status", "Error")
		jsonData.Set("status_code", 0)
		jsonData.Set("message", err.Error())
		sendRes(w, jsonData)
		return
	}

	jsonData.Set("status", "Success")
	jsonData.Set("status_code", "01")
	sendRes(w, jsonData)
}

func PostEditConnParamURL(w http.ResponseWriter, r *http.Request) {
	log.Println("Post Edit")

	// Maximum upload of 10 MB files
	r.ParseMultipartForm(10 << 20)

	var host string = r.FormValue("host")
	var address string = r.FormValue("address")
	var connname string = r.FormValue("name")
	var ConnId string = r.FormValue("connId")

	jsonData := simplejson.New()

	// Get handler for filename, size and headers
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
	}

	var connPath string = Credspath + "/" + ConnId

	if err == nil {
		defer file.Close()

		err := os.RemoveAll(connPath)
		if err != nil {
			log.Fatal(err)
		}

		// This is path which we want to store the file
		if !DirExists(connPath) {
			os.MkdirAll(connPath, 0777)
		}

		f, err := os.OpenFile(connPath+"/privatekey.pem", os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			log.Println("errror in Upload file", err)
			return
		}

		// Copy the file to the destination path
		io.Copy(f, file)
	}

	err = EditConnParameters(host, address, connPath+"/privatekey.pem", connname, ConnId)
	if err != nil {
		log.Println("errror in SQL Add Params", err)
	}

	jsonData.Set("status", "Success")
	jsonData.Set("status_code", "01")

	sendRes(w, jsonData)
}

func PostCreateBranch(w http.ResponseWriter, r *http.Request) {

	log.Println("PostCreateBranch")
	// Accept only POST
	if r.Method != "POST" {
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create empty return JSON
	jsonData := simplejson.New()

	var postBody branchStruct
	// Decode post body
	err := json.NewDecoder(r.Body).Decode(&postBody)
	if err != nil {
		log.Println("Error in UnMarsal")
		jsonData.Set("status", "Error")
		jsonData.Set("status_code", "990")
		jsonData.Set("message", "Cannot Unmarshal Post Body")
		sendRes(w, jsonData)
		return
	}

	log.Println("postBody", postBody.Script)

	postBody.BranchId = genUUID()

	var connPath string = "./GitDir/" + postBody.ConnId + "/" + postBody.BranchId
	// This is path which we want to store the file
	if !DirExists(connPath) {
		os.MkdirAll(connPath, 0777)
	}

	req_code, req_msg := MySQL_PostAddBranch(postBody)
	if req_code == 0 {
		log.Println("Error in Getting Scripts SQL")
		jsonData.Set("status", "Error")
		jsonData.Set("status_code", req_code)
		jsonData.Set("message", req_msg)
		sendRes(w, jsonData)
		return
	}

	log.Println("New Branch Script")
	tempScript := "cd " + connPath + "; " + postBody.Script
	output, err := ExecuteScriptsInDevina(tempScript)
	if err != nil {
		log.Println("Error in Devina Execute", err)
		jsonData.Set("status", "Error")
		jsonData.Set("status_code", req_code)
		jsonData.Set("message", req_msg)
		sendRes(w, jsonData)
		return
	}

	jsonData.Set("status", "Success")
	jsonData.Set("status_code", "01")
	jsonData.Set("output", output)

	sendRes(w, jsonData)
}
