package main

import (
	"encoding/json"
	"flag"
	"os"
)

func initProjects() (resp *Projects, err error) {
	var locProjectMetaFile string

	flag.StringVar(&locProjectMetaFile, "c", "projects/projects.json", "Full path of the metadata file for projects")
	flag.Parse()

	_b, err := os.ReadFile(locProjectMetaFile)
	if err != nil {
		return
	}

	err = json.Unmarshal(_b, &resp)
	return
}
