package main

import (
	"io/ioutil"
	"os"
)

func CreateFile(fileContent string) (bool, error) {

	f, err := os.Create("data.go")
	if err != nil {
		return false, err
	}

	defer f.Close()

	_, err = f.WriteString(fileContent)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ReadFile(Path string) (string, error) {

	content, err := ioutil.ReadFile(Path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
