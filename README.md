GoVersionInfo
==========
[![Build Status](https://travis-ci.org/josephspurrier/goversioninfo.svg)](https://travis-ci.org/josephspurrier/goversioninfo) [![Coverage Status](https://coveralls.io/repos/josephspurrier/goversioninfo/badge.png)](https://coveralls.io/r/josephspurrier/goversioninfo) [![GoDoc](https://godoc.org/github.com/josephspurrier/goversioninfo?status.svg)](https://godoc.org/github.com/josephspurrier/goversioninfo)

Golang Microsoft Version Info and Icon Resource Generator

Package creates a syso file which contains Microsoft Version Information and an optional icon. When you run "go build", Go will embed the version information and icon in the executable. Go will automatically use the syso if it's in the same directory as the main() function. Documentation available on [GoDoc](https://godoc.org/github.com/josephspurrier/goversioninfo).

## Using the Package

Copy versioninfo.json and an icon named icon.ico into your working directory. Fill versioninfo.json with your own information. Use the code below to generate a syso file named resource.syso.

~~~ go
package main

import (
	"fmt"
	"io/ioutil"
	
	"github.com/josephspurrier/goversioninfo"
)

func main() {	
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
~~~

## Major Contributions

Thanks to [Mateusz Czaplinski](https://github.com/akavel/rsrc) for his embedded binary resource package.

Thanks to [Tamás Gulácsi](https://github.com/tgulacsi) for his superb code additions, refactoring, optimization to make this a solid package.