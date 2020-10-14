package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

var owner string = "lakbhat"
var repo string = "SententiaRepo"
var sha string = "dev_ui_editor"
var bearer = "Bearer " + "3b20ea4fe651a1882c8847082478469251f56ffc"
var url string = "https://api.github.com/repos/lakbhat/SententiaRepo/"

type GitResStruct struct {
	Sha         string          `json:"sha,omitempty"`
	NodeId      string          `json:"node_id,omitempty"`
	Commit      commitStruct    `json:"commit,omitempty"`
	Url         string          `json:"url,omitempty"`
	HtmlUrl     string          `json:"html_url,omitempty"`
	CommentsUrl string          `json:"comments_url,omitempty"`
	Author      authorStruct    `json:"author,omitempty"`
	Committer   authorStruct    `json:"committer,omitempty"`
	Parents     []parentsStruct `json:"parents,omitempty"`
}

type commitStruct struct {
	Author       commitAuthorStruct    `json:"author"`
	Committer    commitCommitterStruct `json:"committer"`
	Message      string                `json:"message"`
	Tree         treeStruct            `json:"tree"`
	URL          string                `json:"url"`
	CommentCount int                   `json:"comment_count"`
	Verification verificationStruct    `json:"verification"`
}

type commitAuthorStruct struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

type commitCommitterStruct struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Data  string `json:"date"`
}

type treeStruct struct {
	Sha string `json:"sha"`
	URL string `json:"url"`
}

type verificationStruct struct {
	Verified  bool   `json:"verified"`
	Reason    string `json:"reason"`
	Signature bool   `json:"signature"`
	Payload   bool   `json:"payload"`
}

type authorStruct struct {
	Login               string `json:"login"`
	ID                  int    `json:"id"`
	Node_id             string `json:"node_id"`
	Avatar_url          string `json:"avatar_url"`
	Gravatar_id         string `json:"gravatar_id"`
	Url                 string `json:"url"`
	HtmlURL             string `json:"html_url"`
	Followers_url       string `json:"followers_url"`
	Following_url       string `json:"following_url"`
	Gists_url           string `json:"gists_url"`
	Starred_url         string `json:"starred_url"`
	Subscriptions_url   string `json:"subscriptions_url"`
	Organizations_url   string `json:"organizations_url"`
	Repos_url           string `json:"repos_url"`
	Events_url          string `json:"events_url"`
	Received_events_url string `json:"received_events_url"`
	Typez               string `json:"type"`
	Site_admin          bool   `json:"site_admin"`
}

type parentsStruct struct {
	Sha     string `json:"sha"`
	URL     string `json:"url"`
	HtmlURL string `json:"html_url"`
}

func main() {

	getCommits("dev_rule_engine", "1")
}

func getCommits(branch string, perPage string) {

	var tempUrl string = url + "commits?per_page=" + perPage + "&sha=" + branch
	makeGetRequest(tempUrl)

}

func makeGetRequest(url string) []GitResStruct {
	// make a sample HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)

	// check for response error
	if err != nil {
		log.Fatal(err)
	}

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var gitResCollection []GitResStruct
	err = json.Unmarshal(body, &gitResCollection)

	return gitResCollection
}
