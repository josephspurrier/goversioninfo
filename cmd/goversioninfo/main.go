// Contribution by Tamás Gulácsi

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/josephspurrier/goversioninfo"
)

func main() {
	flagExample := flag.Bool("example", false, "dump out an example versioninfo.json to stdout")

	cfg := goversioninfo.NewCLIConfig()

	flagOut := flag.String("o", cfg.OutputFile, "output file name")
	flagGo := flag.String("gofile", "", "Go output file name (optional)")
	flagPackage := flag.String("gofilepackage", cfg.GoFilePackage, "Go output package name (optional, requires parameter: 'gofile')")
	flagPlatformSpecific := flag.Bool("platform-specific", false, "output i386, amd64, arm and arm64 named resource.syso, ignores -o")
	flagIcon := flag.String("icon", "", "icon file name(s), separated by commas")
	flagApplicationIcon := flag.String("application-icon", "", "icon file for IDI_APPLICATION (window title bar); defaults to -icon if unset")
	flagManifest := flag.String("manifest", "", "manifest file name")
	flagSkipVersion := flag.Bool("skip-versioninfo", false, "skip version info")
	flagPropagateVerStrings := flag.Bool("propagate-ver-strings", false,
		"fill FixedFileInfo version fields using FileVersion and ProductVersion from the StringFileInfo")

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

	goarch := os.Getenv("GOARCH")
	if goarch == "" {
		goarch = runtime.GOARCH
	}
	flag64 := flag.Bool("64", goarch == "amd64" || goarch == "arm64", "generate 64-bit binaries")
	flagarm := flag.Bool("arm", goarch == "arm" || goarch == "arm64", "generate arm binaries")

	flagVerMajor := flag.Int("ver-major", -1, "FileVersion.Major")
	flagVerMinor := flag.Int("ver-minor", -1, "FileVersion.Minor")
	flagVerPatch := flag.Int("ver-patch", -1, "FileVersion.Patch")
	flagVerBuild := flag.Int("ver-build", -1, "FileVersion.Build")

	flagProductVerMajor := flag.Int("product-ver-major", -1, "ProductVersion.Major")
	flagProductVerMinor := flag.Int("product-ver-minor", -1, "ProductVersion.Minor")
	flagProductVerPatch := flag.Int("product-ver-patch", -1, "ProductVersion.Patch")
	flagProductVerBuild := flag.Int("product-ver-build", -1, "ProductVersion.Build")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <versioninfo.json>\n\nPossible flags:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *flagExample {
		io.WriteString(os.Stdout, example)
		return
	}

	cfg.ConfigFile = flag.Arg(0)
	if cfg.ConfigFile == "" {
		cfg.ConfigFile = "versioninfo.json"
	}

	cfg.OutputFile = *flagOut
	cfg.GoFile = *flagGo
	cfg.GoFilePackage = *flagPackage
	cfg.PlatformSpecific = *flagPlatformSpecific
	cfg.IconPath = *flagIcon
	cfg.ApplicationIconPath = *flagApplicationIcon
	cfg.ManifestPath = *flagManifest
	cfg.SkipVersionInfo = *flagSkipVersion
	cfg.PropagateVerStrings = *flagPropagateVerStrings

	cfg.Comment = *flagComment
	cfg.CompanyName = *flagCompany
	cfg.Description = *flagDescription
	cfg.FileVersion = *flagFileVersion
	cfg.InternalName = *flagInternalName
	cfg.Copyright = *flagCopyright
	cfg.Trademark = *flagTrademark
	cfg.OriginalName = *flagOriginalName
	cfg.PrivateBuild = *flagPrivateBuild
	cfg.ProductName = *flagProductName
	cfg.ProductVersion = *flagProductVersion
	cfg.SpecialBuild = *flagSpecialBuild

	cfg.TranslationID = *flagTranslation
	cfg.CharsetID = *flagCharset

	cfg.Is64Bit = *flag64
	cfg.IsARM = *flagarm

	cfg.VerMajor = *flagVerMajor
	cfg.VerMinor = *flagVerMinor
	cfg.VerPatch = *flagVerPatch
	cfg.VerBuild = *flagVerBuild

	cfg.ProductVerMajor = *flagProductVerMajor
	cfg.ProductVerMinor = *flagProductVerMinor
	cfg.ProductVerPatch = *flagProductVerPatch
	cfg.ProductVerBuild = *flagProductVerBuild

	if err := goversioninfo.RunCLI(cfg); err != nil {
		log.Fatal(err)
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
