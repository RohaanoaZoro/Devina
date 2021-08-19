package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

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

type RequiredCommitInfo struct {
	Sha        string `json:"sha"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Date       int    `json:"date"`
	Message    string `json:"message"`
	ParentsSha string `json:"parentsha"`
}

func getCommits(branch string, per_page string) (int, string, []RequiredCommitInfo) {

	var tempUrl string = url + "commits?sha=" + branch
	if per_page != "" {
		tempUrl += "&per_page=" + per_page
	}

	res_code, res_msg, resBody := makeGetRequest(tempUrl)
	if res_code != 1 {
		return res_code, res_msg, nil
	}

	commits := getUsefulInfo(resBody)

	return res_code, res_msg, commits
}

func getCommitsByTime(branch string, perPage string, since string) (int, string, []RequiredCommitInfo) {

	var tempUrl string = url + "commits?sha=" + branch + "&since=" + since
	res_code, res_msg, resBody := makeGetRequest(tempUrl)
	if res_code != 1 {
		return res_code, res_msg, nil
	}

	commits := getUsefulInfo(resBody)

	return res_code, res_msg, commits
}

func getUsefulInfo(resBody []byte) []RequiredCommitInfo {

	var gitResCollection []GitResStruct
	err := json.Unmarshal(resBody, &gitResCollection)
	if err != nil {
		log.Println("Error While UnMarshalling")
	}

	var uji []RequiredCommitInfo

	for i := 0; i < len(gitResCollection); i++ {
		// log.Println(gitResCollection[i])
		gitDate := gitResCollection[i].Commit.Author.Date
		gittime, err := time.Parse(time.RFC3339, gitDate)
		now := time.Now()

		days := int(now.Sub(gittime).Hours() / 24)
		if err != nil {
			log.Println(err)
		}
		var singleuji RequiredCommitInfo = RequiredCommitInfo{Sha: gitResCollection[i].Commit.Tree.Sha, Name: gitResCollection[i].Commit.Author.Name, Email: gitResCollection[i].Commit.Author.Email, Date: days, Message: gitResCollection[i].Commit.Message, ParentsSha: gitResCollection[i].Parents[0].Sha}

		uji = append(uji, singleuji)
	}

	return uji
}

func makeGetRequest(url string) (int, string, []byte) {
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
		return 0, err.Error(), nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
		return 0, err.Error(), nil
	}

	return 1, "Success", body
}

func getLatestCommit(branch string) (int, string, RequiredCommitInfo) {

	var tempUrl string = url + "commits?per_page=1&sha=" + branch
	res_code, res_msg, resBody := makeGetRequest(tempUrl)
	if res_code != 1 {
		return res_code, res_msg, RequiredCommitInfo{}
	}

	var CommitInfo RequiredCommitInfo = getLatestCommitInfo(resBody)

	return res_code, res_msg, CommitInfo
}

func getLatestCommitInfo(resBody []byte) RequiredCommitInfo {

	var gitResCollection []GitResStruct
	err := json.Unmarshal(resBody, &gitResCollection)

	if err != nil {
		log.Println("Error While UnMarshalling")
	}

	i := 0
	var singleuji RequiredCommitInfo = RequiredCommitInfo{Sha: gitResCollection[i].Commit.Tree.Sha, Name: gitResCollection[i].Commit.Author.Name, Email: gitResCollection[i].Commit.Author.Email, Date: 0.0, Message: gitResCollection[i].Commit.Message, ParentsSha: gitResCollection[i].Parents[0].Sha}

	return singleuji
}
