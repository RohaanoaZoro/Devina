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
var bearer = "Bearer " + "e3c16b280dad763681fe1efda88e66e2e0492b97"
var url string = "https://api.github.com/repos/lakbhat/SententiaRepo/"
var since string = ""

type GitResStruct struct {
	Sha         string          `json:"sha ,omitempty"`
	NodeId      string          `json:"node_id,omitempty"`
	Commit      commitStruct    `json:"commit,omitempty"`
	Committer   authorStruct    `json:"committer,omitempty"`
	Message     string          `json:"message,omitempty"`
	Tree        treeStruct      `json:"tree,omitempty"`
	Url         string          `json:"url,omitempty"`
	HtmlUrl     string          `json:"html_url,omitempty"`
	CommentsUrl string          `json:"comments_url,omitempty"`
	Author      authorStruct    `json:"author,omitempty"`
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
	Date  string `json:"date"`
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

type UsefulJsonInfo struct {
	Sha        string `json:"sha"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Date       string `json:"date"`
	Message    string `json:"message"`
	parentsSha string `json:"sha"`
}

func main() {

	// var lastDate string = getLastDateInPreProd("prod_ui_editor")
	// log.Println(lastDate)
	// getCommits("dev_ui_editor", "1", "2021-01-01T12:13:43Z")
	getCommits("preprod_ui_editor", "1", "2021-01-19T16:40:04Z")
	// getCommits("preprod_ui_editor", "1", lastDate)
}

func getCommits(branch string, perPage string, since string) {

	var tempUrl string = url + "commits?sha=" + branch + "&since=" + since
	var resBody []byte = makeGetRequest(tempUrl)
	uji := getUsefulInfo(resBody)
	log.Println(uji)

}

func getLastDateInPreProd(branch string) string {

	var tempUrl string = url + "commits?per_page=1&sha=" + branch
	var resBody []byte = makeGetRequest(tempUrl)
	return getLastDate(resBody)
}

func makeGetRequest(url string) []byte {
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

	return body
}

func getUsefulInfo(resBody []byte) []UsefulJsonInfo {

	var gitResCollection []GitResStruct
	err := json.Unmarshal(resBody, &gitResCollection)

	if err != nil {
		log.Println("Error While UnMarshalling")
	}

	var uji []UsefulJsonInfo

	for i := 0; i < len(gitResCollection); i++ {
		// log.Println(gitResCollection[i])
		log.Println(gitResCollection[i].Commit.Message)
		var singleuji UsefulJsonInfo = UsefulJsonInfo{Sha: gitResCollection[i].Commit.Tree.Sha, Name: gitResCollection[i].Commit.Author.Name, Email: gitResCollection[i].Commit.Author.Email, Date: gitResCollection[i].Commit.Author.Date, Message: gitResCollection[i].Commit.Message, parentsSha: gitResCollection[i].Parents[0].Sha}

		uji = append(uji, singleuji)
	}

	return uji
}

func getLastDate(resBody []byte) string {

	var gitResCollection []GitResStruct
	err := json.Unmarshal(resBody, &gitResCollection)

	if err != nil {
		log.Println("Error While UnMarshalling")
	}

	return gitResCollection[0].Commit.Author.Date
}
