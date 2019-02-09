package goversioninfo

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Example
func Example() {
	// Read the config file
	jsonBytes, err := ioutil.ReadFile("testdata/resource/versioninfo.json")
	if err != nil {
		log.Printf("Error reading versioninfo.json: %v", err)
		os.Exit(1)
	}

	// Create a new container
	vi := &VersionInfo{}

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
	vi.IconPath = "testdata/resource/icon.ico"

	// Create the file
	if err := vi.WriteSyso("resource.syso", "386"); err != nil {
		log.Printf("Error writing syso: %v", err)
		os.Exit(3)
	}
}

func ExampleUseIcon() {
	filename := "cmd"

	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		fmt.Println("Could not parse "+filename+".json", err)
	}

	vi.IconPath = "testdata/resource/icon.ico"

	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	file := "resource.syso"

	vi.WriteSyso(file, "386")

	_, err = ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Could not load "+file, err)
	}
}

func ExampleUseTimestamp() {
	filename := "cmd"

	path, _ := filepath.Abs("./testdata/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		fmt.Println("Could not parse "+filename+".json", err)
	}

	// Write a timestamp even though it is against the spec
	vi.Timestamp = true

	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	file := "resource.syso"

	vi.WriteSyso(file, "386")

	_, err = ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Could not load "+file, err)
	}
}
