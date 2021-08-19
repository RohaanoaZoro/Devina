package main

import (
	"encoding/json"
	"log"
)

type GitBranchStruct struct {
	Name      string             `json:"name"`
	Commit    commitBranchStruct `json:"commit"`
	Protected bool               `json:"protected"`
}

type commitBranchStruct struct {
	Sha string `json:"sha"`
	URL string `json:"url"`
}

func GetGitBranches() (int, string, []string) {

	var tempUrl string = url + "branches"
	res_code, res_msg, resBody := makeGetRequest(tempUrl)
	if res_code != 1 {
		return res_code, res_msg, nil
	}

	var branchArr []GitBranchStruct
	err := json.Unmarshal(resBody, &branchArr)
	if err != nil {
		log.Println("Error While UnMarshalling")
		return res_code, err.Error(), nil
	}

	var finalArr []string
	for i := 0; i < len(branchArr); i++ {
		finalArr = append(finalArr, branchArr[i].Name)
	}

	return 1, "Success", finalArr
}
