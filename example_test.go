// Copyright 2015 Joseph Spurrier
// Author: Joseph Spurrier (http://josephspurrier.com)
// License: http://www.apache.org/licenses/LICENSE-2.0.html

package goversioninfo

import (
	"fmt"
	"io/ioutil"

	"github.com/josephspurrier/goversioninfo"
)

// Example
func Example() {
	logic()
}

// Create the syso file
func logic() {
	// Read the config file
	jsonBytes, err := ioutil.ReadFile("versioninfo.json")
	if err != nil {
		log.Printf("Error reading %q: %v", configFile, err)
		os.Exit(1)
	}

	// Create a new container
	vi := &goversioninfo.VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		log.Printf("Could not parse the .json file: %v", err)
		os.Exit(2)
	}

	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	// Optionally, embed an icon by path
	// If the icon has multiple sizes, all of the sizes will be embedded
	vi.IconPath = "icon.ico"

	// Create the file
	if err := vi.WriteSyso("resource.syso"); err != nil {
		log.Printf("Error writing syso: %v", err)
		os.Exit(3)
	}
}
