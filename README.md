goconsole
==========
[![GoDoc](https://godoc.org/github.com/josephspurrier/goversioninfo?status.svg)](https://godoc.org/github.com/josephspurrier/versioninfo)

Golang Microsoft Version Info and Icon Resource Generator

Package creates syso files with Version Info and Icon that the go build command will pick up and embed in your application. Fill in the versioninfo.json file and then use the code below to generate a .syso file. When you build, Go with AUTOMATICALLY find the file if is is in the same directory as your file with the main() function. Complete documentation available on [GoDoc](https://godoc.org/github.com/josephspurrier/versioninfo).

## Code

```
package main

import (
	"fmt"
	"github.com/josephspurrier/goversioninfo"
	"io/ioutil"
)

func main() {
	// Read the config file
	jsonBytes, err := ioutil.ReadFile("versioninfo.json")
	if err != nil {
		fmt.Println("File Error:", err)
		return
	}

	// Create a new container
	vi := &goversioninfo.VersionInfo{}

	// Parse the config
	if ok := vi.ParseJSON(jsonBytes); ok {
		// Fill the structures with config data
		vi.Build()

		// Write the data to a buffer
		vi.Walk()

		// Create the file
		vi.WriteSyso("resource.syso")
	} else {
		fmt.Println("Could not parse the .json file")
	}
}
```