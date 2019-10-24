package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type RunResult struct {
	ID string

	Result string
	Runs   []string
	Ext    string
}

func getResult(id string) RunResult {
	var result RunResult

	jb, err := ioutil.ReadFile("results/" + id + "/result.json")
	if err != nil {
		log.Println("[Result] Cannot open result file")
		return result
	}

	err = json.Unmarshal(jb, &result)
	if err != nil {
		log.Println("[Result] Cannot open result file")
		return result
	}

	return result
}

func saveResult(id string, result RunResult) {
	jb, _ := json.Marshal(result)

	jsonFile, err := os.Create("results/" + id + "/result.json")
	if err != nil {
		log.Println("[Result] Cannot create result file")
		return
	}
	defer jsonFile.Close()

	jsonFile.Write(jb)
}
