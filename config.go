package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Server struct {
	Host        string   `json:"host"`
	Username    string   `json:"username"`
	Destination string   `json:"destination"`
	PostDeploy  []string `json:"post_deploy"`
}

type Build struct {
	Command string   `json:"command"`
	Output  string   `json:"output"`
	Files   []string `json:"files"`
}

type ConfigFile struct {
	Root   string `json:"root"`
	Server Server `json:"server"`
	Build  Build  `json:"build"`
}

func (s *ConfigFile) load() {
	data, err := os.ReadFile("deploy_config.json")
	if check(err, "No Config file found", false) != true {
		return
	}

	err = json.Unmarshal(data, &s)
	check(err, "Something is wrong with the file, please check", true)

}

func (f *ConfigFile) save() {
	data, err := json.Marshal(f)
	check(err, "Something is wrong with the data that is being put in the file, please check", true)

	fmt.Println(string(data))

	err = os.WriteFile("deploy_config.json", data, 0644)
	check(err, "Something is wrong with writing the data, please raise an issue", true)
}

func Config() *ConfigFile {
	var file ConfigFile
	file.load()
	file.save()
	return &file
}

func check(e error, message string, crucial bool) bool {
	if e != nil {
		fmt.Printf("%s\n", message)
		if crucial {
			os.Exit(1)
		}
		return false
	}
	return true
}
