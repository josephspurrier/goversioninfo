// Copyright 2015 Joseph Spurrier
// Author: Joseph Spurrier (http://josephspurrier.com)
// License: http://www.apache.org/licenses/LICENSE-2.0.html

// Package goversion create syso files with Microsoft Version Information embedded.
package goversioninfo

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/akavel/rsrc/binutil"
	"github.com/akavel/rsrc/coff"
	"github.com/akavel/rsrc/ico"
	"io"
	"io/ioutil"
	"math"
	"os"
	"reflect"
	"strconv"
	//"syscall"
	//"time"
)

// *****************************************************************************
// JSON and Config
// *****************************************************************************

// Parse JSON file to structs
func (vi *VersionInfo) ParseJSON(jsonBytes []byte) bool {
	if err := json.Unmarshal([]byte(jsonBytes), &vi); err != nil {
		fmt.Println("ParseJSON Error:", err)
		return false
	}

	return true
}

// Version Info data container
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

// Translation with langid and charsetid
type Translation struct {
	LangID    string
	CharsetID string
}

// Version with 3 parts
type FileVersion struct {
	Major int
	Minor int
	Patch int
	Build int
}

// File characteristics - leave most of them at the defaults
type FixedFileInfo struct {
	FileVersion    `json:"FileVersion"`
	ProductVersion FileVersion
	FileFlagsMask  string
	FileFlags      string
	FileOS         string
	FileType       string
	FileSubType    string
}

// Translation container
type VarFileInfo struct {
	Translation `json:"Translation"`
}

// Settings you want to change
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
	b := make([]byte, 0)

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
func (v *VersionInfo) Walk() {
	// Create a buffer
	var b bytes.Buffer
	w := binutil.Writer{W: &b}

	// Write to the buffer
	binutil.Walk(v.Structure, func(v reflect.Value, path string) error {
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

	v.Buffer = b
}

// WriteSyso creates a resource file from the version info and optionally an icon
func (v *VersionInfo) WriteSyso(filename string) {

	// Channel for generating IDs
	newid := make(chan uint16)
	go func() {
		for i := uint16(1); ; i++ {
			newid <- i
		}
	}()

	// Create a new RSRC section
	coff := coff.NewRSRC()
	id := <-newid

	// ID 16 is for Version Information
	coff.AddResource(16, id, sizedReader{&v.Buffer})

	// If icon is enabled
	if v.Icon {
		err := addicon(coff, v.IconPath, newid)
		if err != nil {
			//return err
			fmt.Println("Error adding icon:", err)
			return
		}
	}

	coff.Freeze()

	// Write to file
	err := writeCoff(coff, filename)
	if err != nil {
		fmt.Printf("error writing %s coff/.syso: %v\n", err)
		return
	}
}

// WriteHex creates a hex file for debugging version info
func (v *VersionInfo) WriteHex(filename string) {
	ioutil.WriteFile(filename, v.Buffer.Bytes(), 0655)
}

// *****************************************************************************
// Structure Building
// *****************************************************************************

/*
Version Information Structures
http://msdn.microsoft.com/en-us/library/windows/desktop/ff468916.aspx

VersionInfo Names
http://msdn.microsoft.com/en-us/library/windows/desktop/aa381058.aspx#string-name

Translation: LangID
http://msdn.microsoft.com/en-us/library/windows/desktop/aa381058.aspx#langid

Translation: CharsetID
http://msdn.microsoft.com/en-us/library/windows/desktop/aa381058.aspx#charsetid

*/

// Top level version container
type VS_VersionInfo struct {
	WLength      uint16
	WValueLength uint16
	WType        uint16
	SzKey        []byte
	Padding1     []byte
	Value        VS_FixedFileInfo
	Padding2     []byte
	Children     VS_StringFileInfo
	Children2    VS_VarFileInfo
}

// Most of these should be left at the defaults
type VS_FixedFileInfo struct {
	DwSignature        uint32
	DwStrucVersion     uint32
	DwFileVersionMS    uint32
	DwFileVersionLS    uint32
	DwProductVersionMS uint32
	DwProductVersionLS uint32
	DwFileFlagsMask    uint32
	DwFileFlags        uint32
	DwFileOS           uint32
	DwFileType         uint32
	DwFileSubtype      uint32
	DwFileDateMS       uint32
	DwFileDateLS       uint32
}

// Holds multiple collections of keys and values, only allows for 1 collection in
// this package
type VS_StringFileInfo struct {
	WLength      uint16
	WValueLength uint16
	WType        uint16
	SzKey        []byte
	Padding      []byte
	Children     VS_StringTable
}

// Holds a collection of string keys and values
type VS_StringTable struct {
	WLength      uint16
	WValueLength uint16
	WType        uint16
	SzKey        []byte
	Padding      []byte
	Children     []VS_String
}

// Holds the keys and values
type VS_String struct {
	WLength      uint16
	WValueLength uint16
	WType        uint16
	SzKey        []byte
	Padding1     []byte
	Value        []byte
	Padding2     []byte
}

// Holds the translation collection of 1
type VS_VarFileInfo struct {
	WLength      uint16
	WValueLength uint16
	WType        uint16
	SzKey        []byte
	Padding      []byte
	Value        VS_Var
}

// Holds the translation key
type VS_Var struct {
	WLength      uint16
	WValueLength uint16
	WType        uint16
	SzKey        []byte
	Padding      []byte
	Value        uint32
}

func buildString(i int, v reflect.Value) (VS_String, bool, uint16) {
	sValue := string(v.Field(i).Interface().(string))
	sName := v.Type().Field(i).Name

	ss := VS_String{}

	// If the value is set
	if sValue != "" {
		// Create key
		ss.SzKey = buildUnicode(sName, false)
		soFar := len(ss.SzKey) + 6
		ss.Padding1 = padBytes(4 - int(math.Mod(float64(soFar), 4)))
		// Ensure there is at least 4 bytes between the key and value by NOT
		// using this code
		/*if len(ss.Padding1) == 4 {
			ss.Padding1 = []byte{}
		}*/

		// Create value
		ss.Value = buildUnicode(sValue, true)
		soFar += (len(ss.Value) + len(ss.Padding1))
		ss.Padding2 = padBytes(4 - int(math.Mod(float64(soFar), 4)))
		// Eliminate too much spacing
		if len(ss.Padding2) == 4 {
			ss.Padding2 = []byte{}
		}

		// Length of text in words (2 bytes)
		ss.WValueLength = uint16(len(ss.Value) / 2)
		// This is NOT a good way because the copyright symbol counts as 2 letters
		//ss.WValueLength = uint16(len(sValue) + 1)

		// 0 for binary, 1 for text
		ss.WType = 0x01

		// Length of structure
		ss.WLength = uint16(soFar)
		// Don't include the padding in the length, but you must pass it back to
		// the parent to be included
		//ss.WLength = uint16(soFar + len(ss.Padding2))

		return ss, true, uint16(len(ss.Padding2))
	}

	return ss, false, 0
}

func buildStringTable(vi *VersionInfo) (VS_StringTable, uint16) {
	st := VS_StringTable{}

	// Always set to 0
	st.WValueLength = 0x00

	// 0 for binary, 1 for text
	st.WType = 0x01

	// Language identifier and Code page
	st.SzKey = buildUnicode(vi.VarFileInfo.Translation.getTranslationString(), false)
	soFar := len(st.SzKey) + 6
	st.Padding = padBytes(4 - int(math.Mod(float64(soFar), 4)))

	// Loop through the struct fields
	v := reflect.ValueOf(vi.StringFileInfo)
	for i := 0; i < v.NumField(); i++ {
		// If the struct is valid
		if r, ok, extra := buildString(i, v); ok {
			st.Children = append(st.Children, r)
			st.WLength += (r.WLength + extra)
		}
	}

	st.WLength += uint16(soFar)

	return st, uint16(len(st.Padding))
}

func buildStringFileInfo(vi *VersionInfo) (VS_StringFileInfo, uint16) {
	sf := VS_StringFileInfo{}

	// Always set to 0
	sf.WValueLength = 0x00

	// 0 for binary, 1 for text
	sf.WType = 0x01

	sf.SzKey = buildUnicode("StringFileInfo", false)
	soFar := len(sf.SzKey) + 6
	sf.Padding = padBytes(4 - int(math.Mod(float64(soFar), 4)))

	// Allows for more than one string table
	st, extra := buildStringTable(vi)
	sf.Children = st
	sf.WLength += (uint16(soFar) + uint16(len(sf.Padding)) + st.WLength)

	return sf, extra
}

func buildVar(vfi VarFileInfo) VS_Var {
	vs := VS_Var{}
	// Create key
	vs.SzKey = buildUnicode("Translation", false)
	soFar := len(vs.SzKey) + 6
	vs.Padding = padBytes(4 - int(math.Mod(float64(soFar), 4)))

	// Create value
	vs.Value = str2Uint32(vfi.Translation.getTranslation())
	soFar += (4 + len(vs.Padding))

	// Length of text in bytes
	vs.WValueLength = 4

	// 0 for binary, 1 for text
	vs.WType = 0x00

	// Length of structure
	vs.WLength = uint16(soFar)

	return vs
}

func buildVarFileInfo(vfi VarFileInfo) VS_VarFileInfo {
	vf := VS_VarFileInfo{}

	// Always set to 0
	vf.WValueLength = 0x00

	// 0 for binary, 1 for text
	vf.WType = 0x01

	vf.SzKey = buildUnicode("VarFileInfo", false)
	soFar := len(vf.SzKey) + 6
	vf.Padding = padBytes(4 - int(math.Mod(float64(soFar), 4)))

	// Allows for more than one string table
	st := buildVar(vfi)
	vf.Value = st
	vf.WLength += (uint16(soFar) + uint16(len(vf.Padding)) + st.WLength)

	return vf
}

func buildFixedFileInfo(vi *VersionInfo) VS_FixedFileInfo {
	ff := VS_FixedFileInfo{}
	ff.DwSignature = 0xFEEF04BD
	ff.DwStrucVersion = 0x00010000
	ff.DwFileVersionMS = str2Uint32(vi.FixedFileInfo.FileVersion.getVersionHighString())
	ff.DwFileVersionLS = str2Uint32(vi.FixedFileInfo.FileVersion.getVersionLowString())
	ff.DwProductVersionMS = str2Uint32(vi.FixedFileInfo.ProductVersion.getVersionHighString())
	ff.DwProductVersionLS = str2Uint32(vi.FixedFileInfo.ProductVersion.getVersionLowString())
	ff.DwFileFlagsMask = str2Uint32(vi.FixedFileInfo.FileFlagsMask)
	ff.DwFileFlags = str2Uint32(vi.FixedFileInfo.FileFlags)
	ff.DwFileOS = str2Uint32(vi.FixedFileInfo.FileOS)
	ff.DwFileType = str2Uint32(vi.FixedFileInfo.FileType)
	ff.DwFileSubtype = str2Uint32(vi.FixedFileInfo.FileSubType)

	// According to the spec, these should be zero...ugh
	/*if vi.Timestamp {
		now := syscall.NsecToFiletime(time.Now().UnixNano())
		ff.DwFileDateMS = now.HighDateTime
		ff.DwFileDateLS = now.LowDateTime
	}*/

	return ff
}

// Build fills the structs with data from the config file
func (v *VersionInfo) Build() {
	vi := VS_VersionInfo{}

	// 0 for binary, 1 for text
	vi.WType = 0x00

	vi.SzKey = buildUnicode("VS_VERSION_INFO", false)
	soFar := len(vi.SzKey) + 6
	vi.Padding1 = padBytes(4 - int(math.Mod(float64(soFar), 4)))

	vi.Value = buildFixedFileInfo(v)

	// Length of value (always the same)
	vi.WValueLength = 0x34

	// Never needs padding
	vi.Padding2 = []byte{}

	// Build strings
	sf, extraPadding := buildStringFileInfo(v)
	vi.Children = sf

	// Build translation
	vf := buildVarFileInfo(v.VarFileInfo)
	vi.Children2 = vf

	// Calculate the total size
	vi.WLength += (uint16(soFar) + uint16(len(vi.Padding1)) + vi.WValueLength + uint16(len(vi.Padding2)) + vi.Children.WLength + vi.Children2.WLength + extraPadding)

	v.Structure = vi
}

// *****************************************************************************
/*
Code from https://github.com/akavel/rsrc

The MIT License (MIT)

Copyright (c) 2013-2014 The rsrc Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
// *****************************************************************************

const (
	RT_ICON       = coff.RT_ICON
	RT_GROUP_ICON = coff.RT_GROUP_ICON
	RT_MANIFEST   = coff.RT_MANIFEST
)

// on storing icons, see: http://blogs.msdn.com/b/oldnewthing/archive/2012/07/20/10331787.aspx
type GRPICONDIR struct {
	ico.ICONDIR
	Entries []GRPICONDIRENTRY
}

func (group GRPICONDIR) Size() int64 {
	return int64(binary.Size(group.ICONDIR) + len(group.Entries)*binary.Size(group.Entries[0]))
}

type GRPICONDIRENTRY struct {
	ico.IconDirEntryCommon
	Id uint16
}

func addicon(coff *coff.Coff, fname string, newid <-chan uint16) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	//defer f.Close() don't defer, files will be closed by OS when app closes

	icons, err := ico.DecodeHeaders(f)
	if err != nil {
		return err
	}

	if len(icons) > 0 {
		// RT_ICONs
		group := GRPICONDIR{ICONDIR: ico.ICONDIR{
			Reserved: 0, // magic num.
			Type:     1, // magic num.
			Count:    uint16(len(icons)),
		}}
		for _, icon := range icons {
			id := <-newid
			r := io.NewSectionReader(f, int64(icon.ImageOffset), int64(icon.BytesInRes))
			coff.AddResource(RT_ICON, id, r)
			group.Entries = append(group.Entries, GRPICONDIRENTRY{icon.IconDirEntryCommon, id})
		}
		id := <-newid
		coff.AddResource(RT_GROUP_ICON, id, group)
		//fmt.Println("Icon ", fname, " ID: ", id)
	}

	return nil
}

func writeCoff(coff *coff.Coff, fnameout string) error {
	out, err := os.Create(fnameout)
	if err != nil {
		return err
	}
	defer out.Close()
	w := binutil.Writer{W: out}

	// write the resulting file to disk
	binutil.Walk(coff, func(v reflect.Value, path string) error {
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

	if w.Err != nil {
		return fmt.Errorf("Error writing output file: %s", w.Err)
	}

	return nil
}
