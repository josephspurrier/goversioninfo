//go:generate goversioninfo -icon=../resource/icon.ico -manifest=../resource/goversioninfo.exe.manifest -gofile=versioninfo.go

package main

import "fmt"

func main() {
	fmt.Printf("Hello world %v %v %v\n%v", versionInfo.StringFileInfo.ProductName, versionInfo.StringFileInfo.ProductVersion, versionInfo.FixedFileInfo.FileVersion, versionInfo.StringFileInfo.LegalCopyright)
}
