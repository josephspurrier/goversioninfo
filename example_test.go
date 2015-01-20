// Copyright 2015 Joseph Spurrier
// Author: Joseph Spurrier (http://josephspurrier.com)
// License: http://www.apache.org/licenses/LICENSE-2.0.html

package goversioninfo_test

import (
	"fmt"
	"io/ioutil"

	"github.com/josephspurrier/goversioninfo"
)

// Example
func Example() {
	logic()
}

// Read the config file
func logic() {
	// Read the config file
	jsonBytes, err := ioutil.ReadFile("versioninfo.json")
	if err != nil {
		fmt.Println("File Error:", err)
		return
	}

	// Create a new container
	vi := &goversioninfo.VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		fmt.Println("Could not parse the .json file")
	}
	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	// Create the file
	if err := vi.WriteSyso("resource.syso"); err != nil {
		fmt.Println("Could not write resource.syso: %v", err)
	}
}
