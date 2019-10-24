package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Language struct {
	TypeStr   string
	Extension string
	CmdLines  [][]string
}

type LangManager struct {
	Languages []Language
}

var langMgr LangManager

func loadConfig() {
	file, err := ioutil.ReadFile("config.json")

	if err != nil {
		log.Fatalln("[LangMgr] Cannot load config", err)
	}

	err = json.Unmarshal(file, &langMgr)

	if err != nil {
		log.Fatalln("[LangMgr] Cannot load config", err)
	}
}
