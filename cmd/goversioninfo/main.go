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
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/josephspurrier/goversioninfo"
)

func main() {
	flagExample := flag.Bool("example", false, "just dump out an example versioninfo.json to stdout")
	flagOut := flag.String("o", "resource.syso", "output file name")

	flagComment := flag.String("comment", "", "StringFileInfo.Comments")
	flagCompany := flag.String("company", "", "StringFileInfo.CompanyName")
	flagDescription := flag.String("description", "", "StringFileInfo.FileDescription")
	flagFileVersion := flag.String("file-version", "", "StringFileInfo.FileVersion")
	flagInternalName := flag.String("internal-name", "", "StringFileInfo.InternalName")
	flagCopyright := flag.String("copyright", "", "StringFileInfo.LegalCopyright")
	flagTrademark := flag.String("trademark", "", "StringFileInfo.LegalTrademarks")
	flagOriginalName := flag.String("original-name", "", "StringFileInfo.OriginalFilename")
	flagPrivateBuild := flag.String("private-build", "", "StringFileInfo.PrivateBuild")
	flagProductName := flag.String("product-name", "", "StringFileInfo.ProductName")
	flagProductVersion := flag.String("product-version", "", "StringFileInfo.ProductVersion")
	flagSpecialBuild := flag.String("special-build", "", "StringFileInfo.SpecialBuild")

	flagTranslation := flag.Int("translation", 0, "translation ID")
	flagCharset := flag.Int("charset", 0, "charset ID")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <versioninfo.json>\n\nPossible flags:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if *flagExample {
		io.WriteString(os.Stdout, example)
		return
	}

	configFile := flag.Arg(0)
	if configFile == "" {
		configFile = "versioninfo.json"
	}
	var err error
	var input = io.ReadCloser(os.Stdin)
	if configFile != "-" {
		if input, err = os.Open(configFile); err != nil {
			log.Printf("Cannot open %q: %v", configFile, err)
			os.Exit(1)
		}
	}

	// Read the config file
	jsonBytes, err := ioutil.ReadAll(input)
	input.Close()
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

	// Override from flags
	if *flagComment != "" {
		vi.StringFileInfo.Comments = *flagComment
	}
	if *flagCompany != "" {
		vi.StringFileInfo.CompanyName = *flagCompany
	}
	if *flagDescription != "" {
		vi.StringFileInfo.FileDescription = *flagDescription
	}
	if *flagFileVersion != "" {
		vi.StringFileInfo.FileVersion = *flagFileVersion
	}
	if *flagInternalName != "" {
		vi.StringFileInfo.InternalName = *flagInternalName
	}
	if *flagCopyright != "" {
		vi.StringFileInfo.LegalCopyright = *flagCopyright
	}
	if *flagTrademark != "" {
		vi.StringFileInfo.LegalTrademarks = *flagTrademark
	}
	if *flagOriginalName != "" {
		vi.StringFileInfo.OriginalFilename = *flagOriginalName
	}
	if *flagPrivateBuild != "" {
		vi.StringFileInfo.PrivateBuild = *flagPrivateBuild
	}
	if *flagProductName != "" {
		vi.StringFileInfo.ProductName = *flagProductName
	}
	if *flagProductVersion != "" {
		vi.StringFileInfo.ProductVersion = *flagProductVersion
	}
	if *flagSpecialBuild != "" {
		vi.StringFileInfo.SpecialBuild = *flagSpecialBuild
	}

	if *flagTranslation > 0 {
		vi.VarFileInfo.Translation.LangID = goversioninfo.LangID(*flagTranslation)
	}
	if *flagCharset > 0 {
		vi.VarFileInfo.Translation.CharsetID = goversioninfo.CharsetID(*flagCharset)
	}

	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	// Create the file
	if err := vi.WriteSyso(*flagOut); err != nil {
		log.Printf("Error writing syso: %v", err)
		os.Exit(3)
	}
}

const example = `{
	"FixedFileInfo": {
		"FileVersion": {
			"Major": 6,
			"Minor": 3,
			"Patch": 9600,
			"Build": 17284
		},
		"ProductVersion": {
			"Major": 6,
			"Minor": 3,
			"Patch": 9600,
			"Build": 17284
		},
		"FileFlagsMask": "3f",
		"FileFlags ": "00",
		"FileOS": "040004",
		"FileType": "01",
		"FileSubType": "00"
	},
	"StringFileInfo": {
		"Comments": "",
		"CompanyName": "Joseph Spurrier Ltd.",
		"FileDescription": "",
		"FileVersion": "6.3.9600.17284 (aaa.140822-1915)",
		"InternalName": "goversioninfo",
		"LegalCopyright": "© Joseph Spurrier. Licensed under the Apache License, Version 2.0",
		"LegalTrademarks": "",
		"OriginalFilename": "goversioninfo",
		"PrivateBuild": "",
		"ProductName": "Go Version Info",
		"ProductVersion": "6.3.9600.17284",
		"SpecialBuild": ""
	},
	"VarFileInfo": {
		"Translation": {
			"LangID": "0409",
			"CharsetID": "04B0"
		}
	}
}`
