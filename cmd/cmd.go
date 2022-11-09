package cmd

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/josephspurrier/goversioninfo"
)

type Arguments struct {
	FlagExample          *bool
	FlagOut              *string
	FlagGo               *string
	FlagPackage          *string
	FlagPlatformSpecific *bool
	FlagIcon             *string
	FlagManifest         *string
	FlagSkipVersion      *bool

	FlagComment        *string
	FlagCompany        *string
	FlagDescription    *string
	FlagFileVersion    *string
	FlagInternalName   *string
	FlagCopyright      *string
	FlagTrademark      *string
	FlagOriginalName   *string
	FlagPrivateBuild   *string
	FlagProductName    *string
	FlagProductVersion *string
	FlagSpecialBuild   *string

	FlagTranslation *int
	FlagCharset     *int

	Flag64  *bool
	Flagarm *bool

	FlagVerMajor *int
	FlagVerMinor *int
	FlagVerPatch *int
	FlagVerBuild *int

	FlagProductVerMajor *int
	FlagProductVerMinor *int
	FlagProductVerPatch *int
	FlagProductVerBuild *int
}

func Cmd(arg Arguments) {
	if *arg.FlagExample {
		io.WriteString(os.Stdout, example)
		return
	}

	configFile := flag.Arg(0)
	if configFile == "" {
		configFile = "versioninfo.json"
	}

	// Create a new container.
	vi := &goversioninfo.VersionInfo{}

	if !*arg.FlagSkipVersion {
		var err error
		var input = io.ReadCloser(os.Stdin)
		if configFile != "-" {
			if input, err = os.Open(configFile); err != nil {
				log.Printf("Cannot open %q: %v", configFile, err)
				os.Exit(1)
			}
		}

		// Read the config file.
		jsonBytes, err := ioutil.ReadAll(input)
		input.Close()
		if err != nil {
			log.Printf("Error reading %q: %v", configFile, err)
			os.Exit(1)
		}

		// Parse the config.
		if err := vi.ParseJSON(jsonBytes); err != nil {
			log.Printf("Could not parse the .json file: %v", err)
			os.Exit(2)
		}

	}

	// Override from flags.
	if *arg.FlagIcon != "" {
		vi.IconPath = *arg.FlagIcon
	}
	if *arg.FlagManifest != "" {
		vi.ManifestPath = *arg.FlagManifest
	}
	if *arg.FlagComment != "" {
		vi.StringFileInfo.Comments = *arg.FlagComment
	}
	if *arg.FlagCompany != "" {
		vi.StringFileInfo.CompanyName = *arg.FlagCompany
	}
	if *arg.FlagDescription != "" {
		vi.StringFileInfo.FileDescription = *arg.FlagDescription
	}
	if *arg.FlagFileVersion != "" {
		vi.StringFileInfo.FileVersion = *arg.FlagFileVersion
	}
	if *arg.FlagInternalName != "" {
		vi.StringFileInfo.InternalName = *arg.FlagInternalName
	}
	if *arg.FlagCopyright != "" {
		vi.StringFileInfo.LegalCopyright = *arg.FlagCopyright
	}
	if *arg.FlagTrademark != "" {
		vi.StringFileInfo.LegalTrademarks = *arg.FlagTrademark
	}
	if *arg.FlagOriginalName != "" {
		vi.StringFileInfo.OriginalFilename = *arg.FlagOriginalName
	}
	if *arg.FlagPrivateBuild != "" {
		vi.StringFileInfo.PrivateBuild = *arg.FlagPrivateBuild
	}
	if *arg.FlagProductName != "" {
		vi.StringFileInfo.ProductName = *arg.FlagProductName
	}
	if *arg.FlagProductVersion != "" {
		vi.StringFileInfo.ProductVersion = *arg.FlagProductVersion
	}
	if *arg.FlagSpecialBuild != "" {
		vi.StringFileInfo.SpecialBuild = *arg.FlagSpecialBuild
	}

	if *arg.FlagTranslation > 0 {
		vi.VarFileInfo.Translation.LangID = goversioninfo.LangID(*arg.FlagTranslation)
	}
	if *arg.FlagCharset > 0 {
		vi.VarFileInfo.Translation.CharsetID = goversioninfo.CharsetID(*arg.FlagCharset)
	}

	// File Version flags.
	if *arg.FlagVerMajor >= 0 {
		vi.FixedFileInfo.FileVersion.Major = *arg.FlagVerMajor
	}
	if *arg.FlagVerMinor >= 0 {
		vi.FixedFileInfo.FileVersion.Minor = *arg.FlagVerMinor
	}
	if *arg.FlagVerPatch >= 0 {
		vi.FixedFileInfo.FileVersion.Patch = *arg.FlagVerPatch
	}
	if *arg.FlagVerBuild >= 0 {
		vi.FixedFileInfo.FileVersion.Build = *arg.FlagVerBuild
	}

	// Product Version flags.
	if *arg.FlagProductVerMajor >= 0 {
		vi.FixedFileInfo.ProductVersion.Major = *arg.FlagProductVerMajor
	}
	if *arg.FlagProductVerMinor >= 0 {
		vi.FixedFileInfo.ProductVersion.Minor = *arg.FlagProductVerMinor
	}
	if *arg.FlagProductVerPatch >= 0 {
		vi.FixedFileInfo.ProductVersion.Patch = *arg.FlagProductVerPatch
	}
	if *arg.FlagProductVerBuild >= 0 {
		vi.FixedFileInfo.ProductVersion.Build = *arg.FlagProductVerBuild
	}

	// Fill the structures with config data.
	vi.Build()

	// Write the data to a buffer.
	vi.Walk()

	// If the flag is set, then generate the optional Go file.
	if *arg.FlagGo != "" {
		vi.WriteGo(*arg.FlagGo, *arg.FlagPackage)
	}

	// List of the architectures to output.
	var archs []string

	// If platform specific, then output all the architectures for Windows.
	if arg.FlagPlatformSpecific != nil && *arg.FlagPlatformSpecific {
		archs = []string{"386", "amd64", "arm", "arm64"}
	} else {
		// Set the architecture, defaulted to 386(32-bit x86).
		archs = []string{"386"} // 386(32-bit x86)
		if arg.Flagarm != nil && *arg.Flagarm {
			if arg.Flag64 != nil && *arg.Flag64 {
				archs = []string{"arm64"} // arm64(64-bit arm64)
			} else {
				archs = []string{"arm"} // arm(32-bit arm)
			}
		} else {
			if arg.Flag64 != nil && *arg.Flag64 {
				archs = []string{"amd64"} // amd64(64-bit x86_64)
			}
		}

	}

	// Loop through each artchitecture.
	for _, item := range archs {
		// Create the file using the -o argument.
		fileout := *arg.FlagOut

		// If the number of architectures is greater than one, then don't use
		// the -o argument.
		if len(archs) > 1 {
			fileout = fmt.Sprintf("resource_windows_%v.syso", item)
		}

		// Create the file for the specified architecture.
		if err := vi.WriteSyso(fileout, item); err != nil {
			log.Printf("Error writing syso: %v", err)
			os.Exit(3)
		}
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
		"CompanyName": "Company, Inc.",
		"FileDescription": "",
		"FileVersion": "6.3.9600.17284 (aaa.140822-1915)",
		"InternalName": "goversioninfo",
		"LegalCopyright": "© Author. Licensed under MIT.",
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
