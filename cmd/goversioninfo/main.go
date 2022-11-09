// Contribution by Tamás Gulácsi

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/josephspurrier/goversioninfo/cmd"
)

func main() {
	arg := cmd.Arguments{
		FlagExample:          flag.Bool("example", false, "dump out an example versioninfo.json to stdout"),
		FlagOut:              flag.String("o", "resource.syso", "output file name"),
		FlagGo:               flag.String("gofile", "", "Go output file name (optional),"),
		FlagPackage:          flag.String("gofilepackage", "main", "Go output package name (optional, requires parameter: 'gofile'),"),
		FlagPlatformSpecific: flag.Bool("platform-specific", false, "output i386, amd64, arm and arm64 named resource.syso, ignores -o"),
		FlagIcon:             flag.String("icon", "", "icon file name"),
		FlagManifest:         flag.String("manifest", "", "manifest file name"),
		FlagSkipVersion:      flag.Bool("skip-versioninfo", false, "skip version info"),

		FlagComment:        flag.String("comment", "", "StringFileInfo.Comments"),
		FlagCompany:        flag.String("company", "", "StringFileInfo.CompanyName"),
		FlagDescription:    flag.String("description", "", "StringFileInfo.FileDescription"),
		FlagFileVersion:    flag.String("file-version", "", "StringFileInfo.FileVersion"),
		FlagInternalName:   flag.String("internal-name", "", "StringFileInfo.InternalName"),
		FlagCopyright:      flag.String("copyright", "", "StringFileInfo.LegalCopyright"),
		FlagTrademark:      flag.String("trademark", "", "StringFileInfo.LegalTrademarks"),
		FlagOriginalName:   flag.String("original-name", "", "StringFileInfo.OriginalFilename"),
		FlagPrivateBuild:   flag.String("private-build", "", "StringFileInfo.PrivateBuild"),
		FlagProductName:    flag.String("product-name", "", "StringFileInfo.ProductName"),
		FlagProductVersion: flag.String("product-version", "", "StringFileInfo.ProductVersion"),
		FlagSpecialBuild:   flag.String("special-build", "", "StringFileInfo.SpecialBuild"),

		FlagTranslation: flag.Int("translation", 0, "translation ID"),
		FlagCharset:     flag.Int("charset", 0, "charset ID"),

		Flag64:  flag.Bool("64", false, "generate 64-bit binaries"),
		Flagarm: flag.Bool("arm", false, "generate arm binaries"),

		FlagVerMajor: flag.Int("ver-major", -1, "FileVersion.Major"),
		FlagVerMinor: flag.Int("ver-minor", -1, "FileVersion.Minor"),
		FlagVerPatch: flag.Int("ver-patch", -1, "FileVersion.Patch"),
		FlagVerBuild: flag.Int("ver-build", -1, "FileVersion.Build"),

		FlagProductVerMajor: flag.Int("product-ver-major", -1, "ProductVersion.Major"),
		FlagProductVerMinor: flag.Int("product-ver-minor", -1, "ProductVersion.Minor"),
		FlagProductVerPatch: flag.Int("product-ver-patch", -1, "ProductVersion.Patch"),
		FlagProductVerBuild: flag.Int("product-ver-build", -1, "ProductVersion.Build"),
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <versioninfo.json>\n\nPossible flags:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	cmd.Cmd(arg)
}
