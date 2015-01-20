// Copyright 2015 Joseph Spurrier
// Author: Joseph Spurrier (http://josephspurrier.com)
// License: http://www.apache.org/licenses/LICENSE-2.0.html

// Package goversioninfo create syso files with Microsoft Version Information embedded.
package goversioninfo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"

	"github.com/akavel/rsrc/binutil"
	"github.com/akavel/rsrc/coff"
)

// *****************************************************************************
// JSON and Config
// *****************************************************************************

// ParseJSON parses the given bytes as a VersionInfo JSON.
func (vi *VersionInfo) ParseJSON(jsonBytes []byte) error {
	return json.Unmarshal([]byte(jsonBytes), &vi)
}

// VersionInfo data container
type VersionInfo struct {
	FixedFileInfo  `json:"FixedFileInfo"`
	StringFileInfo `json:"StringFileInfo"`
	VarFileInfo    `json:"VarFileInfo"`
	Timestamp      bool
	Buffer         bytes.Buffer
	Structure      VS_VersionInfo
	Icon           bool
	IconPath       string
}

// Translation with langid and charsetid.
type Translation struct {
	LangID    string
	CharsetID string
}

// FileVersion with 3 parts.
type FileVersion struct {
	Major int
	Minor int
	Patch int
	Build int
}

// FixedFileInfo contains file characteristics - leave most of them at the defaults.
type FixedFileInfo struct {
	FileVersion    `json:"FileVersion"`
	ProductVersion FileVersion
	FileFlagsMask  string
	FileFlags      string
	FileOS         string
	FileType       string
	FileSubType    string
}

// VarFileInfo is the translation container.
type VarFileInfo struct {
	Translation `json:"Translation"`
}

// StringFileInfo is what you want to change.
type StringFileInfo struct {
	Comments         string
	CompanyName      string
	FileDescription  string
	FileVersion      string
	InternalName     string
	LegalCopyright   string
	LegalTrademarks  string
	OriginalFilename string
	PrivateBuild     string
	ProductName      string
	ProductVersion   string
	SpecialBuild     string
}

// *****************************************************************************
// Helpers
// *****************************************************************************

type sizedReader struct {
	*bytes.Buffer
}

func (s sizedReader) Size() int64 {
	return int64(s.Buffer.Len())
}

func str2Uint32(s string) uint32 {
	u, err := strconv.ParseUint(s, 16, 32)
	if s == "" {
		return 0
	} else if err != nil {
		fmt.Println("Error parsing uint32:", s, err)
		return 0
	}

	return uint32(u)
}

func buildUnicode(s string, zeroTerminate bool) []byte {
	b := make([]byte, 0, len(s)*2+2)

	for _, x := range s {
		b = append(b, byte(x))
		b = append(b, 0x00)
	}

	if zeroTerminate {
		b = append(b, 0x00)
		b = append(b, 0x00)
	}

	return b
}

func padBytes(i int) []byte {
	return make([]byte, i)
}

func (f FileVersion) getVersionHighString() string {
	return fmt.Sprintf("%04x%04x", f.Major, f.Minor)
}

func (f FileVersion) getVersionLowString() string {
	return fmt.Sprintf("%04x%04x", f.Patch, f.Build)
}

// GetVersionString returns a string representation of the version
func (f FileVersion) GetVersionString() string {
	return fmt.Sprintf("%d.%d.%d.%d", f.Major, f.Minor, f.Patch, f.Build)
}

func (t Translation) getTranslationString() string {
	return fmt.Sprintf("%04X%04X", str2Uint32(t.LangID), str2Uint32(t.CharsetID))
}

func (t Translation) getTranslation() string {
	return fmt.Sprintf("%04x%04x", str2Uint32(t.CharsetID), str2Uint32(t.LangID))
}

// *****************************************************************************
// IO Methods
// *****************************************************************************

// Walk fills the data buffer with hexidecimal data from the structs
func (vi *VersionInfo) Walk() {
	// Create a buffer
	var b bytes.Buffer
	w := binutil.Writer{W: &b}

	// Write to the buffer
	binutil.Walk(vi.Structure, func(v reflect.Value, path string) error {
		if binutil.Plain(v.Kind()) {
			w.WriteLE(v.Interface())
			return nil
		}
		vv, ok := v.Interface().(binutil.SizedReader)
		if ok {
			w.WriteFromSized(vv)
			return binutil.WALK_SKIP
		}
		return nil
	})

	vi.Buffer = b
}

// WriteSyso creates a resource file from the version info and optionally an icon
func (vi *VersionInfo) WriteSyso(filename string) error {

	// Channel for generating IDs
	newID := make(chan uint16)
	go func() {
		for i := uint16(1); ; i++ {
			newID <- i
		}
	}()

	// Create a new RSRC section
	coff := coff.NewRSRC()
	//id := <-newID

	// ID 16 is for Version Information
	coff.AddResource(16, 1, sizedReader{&vi.Buffer})

	// If icon is enabled
	if vi.Icon {
		if err := addIcon(coff, vi.IconPath, newID); err != nil {
			return err
		}
	}

	coff.Freeze()

	// Write to file
	return writeCoff(coff, filename)
}

// WriteHex creates a hex file for debugging version info
func (vi *VersionInfo) WriteHex(filename string) error {
	return ioutil.WriteFile(filename, vi.Buffer.Bytes(), 0655)
}
