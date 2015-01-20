// Copyright 2015 Tamás Gulácsi
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/josephspurrier/goversioninfo"
)

func main() {
	flagOut := flag.String("o", "resource.syso", "output file name")
	flag.Parse()

	configFile := flag.Arg(0)
	if configFile == "" {
		configFile = "versioninfo.json"
	}

	// Read the config file
	jsonBytes, err := ioutil.ReadFile(configFile)
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

	// Create the file
	vi.WriteSyso(*flagOut)
}
