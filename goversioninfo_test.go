// Copyright 2015 Joseph Spurrier
// Author: Joseph Spurrier (http://josephspurrier.com)
// License: http://www.apache.org/licenses/LICENSE-2.0.html

// Package main_test performs testing on main package
package goversioninfo_test

import (
	"bytes"
	"fmt"
	"github.com/josephspurrier/goversioninfo"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// *****************************************************************************
// Logic Testing
// *****************************************************************************

func TestFile1(t *testing.T) {
	testFile(t, "cmd")
	testFile(t, "explorer")
	testFile(t, "control")
}

func testFile(t *testing.T, filename string) {
	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &goversioninfo.VersionInfo{}

	// Parse the config
	if ok := vi.ParseJSON(jsonBytes); ok {
		// Fill the structures with config data
		vi.Build()

		// Write the data to a buffer
		vi.Walk()

		path2, _ := filepath.Abs("./tests/" + filename + ".hex")

		expected, err := ioutil.ReadFile(path2)
		if err != nil {
			t.Error("Could not load "+filename+".hex", err)
		}

		if !bytes.Equal(vi.Buffer.Bytes(), expected) {
			t.Error("Data does not match " + filename + ".hex")
		}
	} else {
		t.Error("Could not parse "+filename+".json", err)
	}
}

func TestWrite(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &goversioninfo.VersionInfo{}

	// Parse the config
	if ok := vi.ParseJSON(jsonBytes); ok {
		// Fill the structures with config data
		vi.Build()

		// Write the data to a buffer
		vi.Walk()

		file := "resource.syso"

		vi.WriteSyso(file)

		_, err = ioutil.ReadFile(file)
		if err != nil {
			t.Error("Could not load "+file, err)
		} else {
			os.Remove(file)
		}
	} else {
		t.Error("Could not parse "+filename+".json", err)
	}
}

func TestMalformedJSON(t *testing.T) {
	filename := "bad"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &goversioninfo.VersionInfo{}

	// Parse the config and return false
	if ok := vi.ParseJSON(jsonBytes); ok {
		t.Error("Application was supposed to return false", err)
	}
}

func TestIcon(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &goversioninfo.VersionInfo{}

	// Parse the config
	if ok := vi.ParseJSON(jsonBytes); ok {

		vi.Icon = true

		vi.IconPath = "icon.ico"

		// Fill the structures with config data
		vi.Build()

		// Write the data to a buffer
		vi.Walk()

		file := "resource.syso"

		vi.WriteSyso(file)

		_, err = ioutil.ReadFile(file)
		if err != nil {
			t.Error("Could not load "+file, err)
		} else {
			os.Remove(file)
		}
	} else {
		t.Error("Could not parse "+filename+".json", err)
	}
}

func TestBadIcon(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &goversioninfo.VersionInfo{}

	// Parse the config
	if ok := vi.ParseJSON(jsonBytes); ok {

		vi.Icon = true
		vi.IconPath = "icon2.ico"

		// Fill the structures with config data
		vi.Build()

		// Write the data to a buffer
		vi.Walk()

		file := "resource.syso"

		vi.WriteSyso(file)

		_, err = ioutil.ReadFile(file)
		if err != nil {
			os.Remove(file)
		} else {
			t.Error("File should not exist "+file, err)
		}
	} else {
		t.Error("Could not parse "+filename+".json", err)
	}
}

func TestTimestamp(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &goversioninfo.VersionInfo{}

	// Parse the config
	if ok := vi.ParseJSON(jsonBytes); ok {

		vi.Timestamp = true

		// Fill the structures with config data
		vi.Build()

		// Write the data to a buffer
		vi.Walk()

		file := "resource.syso"

		vi.WriteSyso(file)

		_, err = ioutil.ReadFile(file)
		if err != nil {
			t.Error("Could not load "+file, err)
		} else {
			os.Remove(file)
		}
	} else {
		t.Error("Could not parse "+filename+".json", err)
	}
}

func TestVersionString(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &goversioninfo.VersionInfo{}

	// Parse the config
	if ok := vi.ParseJSON(jsonBytes); ok {
		if vi.FixedFileInfo.GetVersionString() != "6.3.9600.16384" {
			t.Errorf("Version String does not match: %v", vi.FixedFileInfo.GetVersionString())
		}
	} else {
		t.Error("Could not parse "+filename+".json", err)
	}
}

func TestWriteHex(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &goversioninfo.VersionInfo{}

	// Parse the config
	if ok := vi.ParseJSON(jsonBytes); ok {
		// Fill the structures with config data
		vi.Build()

		// Write the data to a buffer
		vi.Walk()

		file := "resource.syso"

		vi.WriteHex(file)

		_, err = ioutil.ReadFile(file)
		if err != nil {
			t.Error("Could not load "+file, err)
		} else {
			os.Remove(file)
		}
	} else {
		t.Error("Could not parse "+filename+".json", err)
	}
}

// *****************************************************************************
// Examples
// *****************************************************************************

func ExampleUseIcon() {
	filename := "cmd"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &goversioninfo.VersionInfo{}

	// Parse the config
	if ok := vi.ParseJSON(jsonBytes); ok {

		vi.Icon = true

		vi.IconPath = "icon.ico"

		// Fill the structures with config data
		vi.Build()

		// Write the data to a buffer
		vi.Walk()

		file := "resource.syso"

		vi.WriteSyso(file)

		_, err = ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("Could not load "+file, err)
		}
	} else {
		fmt.Println("Could not parse "+filename+".json", err)
	}
}

func ExampleUseTimestamp() {
	filename := "cmd"

	path, _ := filepath.Abs("./tests/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Could not load "+filename+".json", err)
	}

	// Create a new container
	vi := &goversioninfo.VersionInfo{}

	// Parse the config
	if ok := vi.ParseJSON(jsonBytes); ok {

		// Write a timestamp even though it is against the spec
		vi.Timestamp = true

		// Fill the structures with config data
		vi.Build()

		// Write the data to a buffer
		vi.Walk()

		file := "resource.syso"

		vi.WriteSyso(file)

		_, err = ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("Could not load "+file, err)
		}
	} else {
		fmt.Println("Could not parse "+filename+".json", err)
	}
}
