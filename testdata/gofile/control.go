// Auto-generated file by goversioninfo. Do not edit.
package main

import (
	"encoding/json"

	"github.com/josephspurrier/goversioninfo"
)

func unmarshalGoVersionInfo(b []byte) goversioninfo.VersionInfo {
	vi := goversioninfo.VersionInfo{}
	json.Unmarshal(b, &vi)
	return vi
}

var versionInfo = unmarshalGoVersionInfo([]byte(`{
	"FixedFileInfo":{
		"FileVersion": {
			"Major": 6,
			"Minor": 3,
			"Patch": 9600,
			"Build": 16384
		},
		"ProductVersion": {
			"Major": 6,
			"Minor": 3,
			"Patch": 9600,
			"Build": 16384
		},
		"FileFlagsMask": "3f",
		"FileFlags": "",
		"FileOS": "040004",
		"FileType": "01",
		"FileSubType": "00"
	},
	"StringFileInfo":{
		"Comments": "",
		"CompanyName": "Microsoft Corporation",
		"FileDescription": "Windows Control Panel",
		"FileVersion": "6.3.9600.16384 (winblue_rtm.130821-1623)",
		"InternalName": "Control",
		"LegalCopyright": "© Microsoft Corporation. All rights reserved.",
		"LegalTrademarks": "",
		"OriginalFilename": "CONTROL.EXE",
		"PrivateBuild": "",
		"ProductName": "Microsoft® Windows® Operating System",
		"ProductVersion": "6.3.9600.16384",
		"SpecialBuild": ""
	},
	"VarFileInfo":{
		"Translation": {
			"LangID": 1033,
			"CharsetID": 1200
		}
	}
}`))
