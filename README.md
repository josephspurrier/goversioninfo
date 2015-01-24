GoVersionInfo
==========
[![Build Status](https://travis-ci.org/josephspurrier/goversioninfo.svg)](https://travis-ci.org/josephspurrier/goversioninfo) [![Coverage Status](https://coveralls.io/repos/josephspurrier/goversioninfo/badge.png)](https://coveralls.io/r/josephspurrier/goversioninfo) [![GoDoc](https://godoc.org/github.com/josephspurrier/goversioninfo?status.svg)](https://godoc.org/github.com/josephspurrier/goversioninfo)

Microsoft Version Info and Icon Resource Generator for the Go Language

Package creates a syso file which contains Microsoft Version Information and an optional icon. When you run "go build", Go will embed the version information and icon in the executable. Go will automatically use the syso file if it's in the same directory as the main() function.

## Usage

To install, run the following command:

~~~
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
~~~

Copy versioninfo.json into your working directory and then modify the file with your own settings.

Then to utilize the wonderful "go generate" command, add a similar text to the top of your Go source code:
~~~ go
//go:generate goversioninfo -icon=icon.ico
~~~

Run "go generate" before each "go build" and goversioninfo.exe will create a file called resource.syso in the same directory as the Go source code.

## Command-Line Flags

A complete list of the flags you can pass to goversioninfo are below:

~~~
  -charset=0: charset ID
  -comment="": StringFileInfo.Comments
  -company="": StringFileInfo.CompanyName
  -copyright="": StringFileInfo.LegalCopyright
  -description="": StringFileInfo.FileDescription
  -example=false: just dump out an example versioninfo.json to stdout
  -file-version="": StringFileInfo.FileVersion
  -icon="": icon file name
  -internal-name="": StringFileInfo.InternalName
  -o="resource.syso": output file name
  -original-name="": StringFileInfo.OriginalFilename
  -private-build="": StringFileInfo.PrivateBuild
  -product-name="": StringFileInfo.ProductName
  -product-version="": StringFileInfo.ProductVersion
  -special-build="": StringFileInfo.SpecialBuild
  -trademark="": StringFileInfo.LegalTrademarks
  -translation=0: translation ID
  -ver-build=-1: FileVersion.Build
  -ver-major=-1: FileVersion.Major
  -ver-minor=-1: FileVersion.Minor
  -ver-patch=-1: FileVersion.Patch
~~~

You can look over the Microsoft Resource Information: [VERSIONINFO resource](https://msdn.microsoft.com/en-us/library/windows/desktop/aa381058(v=vs.85).aspx)

You can look through the Microsoft Version Information structures: [Version Information Structures](https://msdn.microsoft.com/en-us/library/windows/desktop/ff468916(v=vs.85).aspx)

## Major Contributions

Thanks to [Tamás Gulácsi](https://github.com/tgulacsi) for his superb code additions, refactoring, optimization to make this a solid package.

Thanks to [Mateusz Czaplinski](https://github.com/akavel/rsrc) for his embedded binary resource package.