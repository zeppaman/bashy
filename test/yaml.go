package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Parameter struct {
	Name string
	Desc string
	Type string
}

type Script struct {
	Cmds   []string // map[string][]string
	Cmd    string
	Name   string
	Params []Parameter
}

func main() {
	filename, _ := filepath.Abs("samples/bash.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config Script

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Value: %#v\n", prettyPrint(config))
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
