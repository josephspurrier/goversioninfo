package goversioninfo

import (
	"bytes"
	"debug/pe"
	"encoding/binary"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/akavel/rsrc/coff"
	"github.com/stretchr/testify/assert"
)

// *****************************************************************************
// Logic Testing
// *****************************************************************************

func TestFile1(t *testing.T) {
	testFile(t, "cmd")
	testFile(t, "explorer")
	testFile(t, "control")
	testFile(t, "simple")
}

func testFile(t *testing.T, filename string) {
	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := os.ReadFile(path)
	assert.NoError(t, err)

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Error("Could not parse "+filename+".json", err)
	}
	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	path2, _ := filepath.Abs("./testdata/hex/" + filename + ".hex")

	// This is for easily exporting results when the algorithm improves
	/*path3, _ := filepath.Abs("./testdata/" + filename + ".out")
	os.WriteFile(path3, vi.Buffer.Bytes(), 0655)*/

	expected, err := os.ReadFile(path2)
	assert.NoError(t, err)

	if !bytes.Equal(vi.Buffer.Bytes(), expected) {
		t.Error("Data does not match " + filename + ".hex")
	}

	// Test the Go file generation.
	tmpdir, err := os.MkdirTemp("", "generate_go")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	path4 := filepath.Join(tmpdir, filename+".go")
	err = vi.WriteGo(path4, "")
	assert.NoError(t, err)

	gen, err := os.ReadFile(path4)
	assert.NoError(t, err)

	path5, _ := filepath.Abs("./testdata/gofile/" + filename + ".go")
	expected5, err := os.ReadFile(path5)
	if err != nil {
		t.Error("Could not load "+path5, err)
	}

	// Handle newlines.
	if runtime.GOOS == "windows" {
		expected5 = []byte(strings.ReplaceAll(string(expected5), "\r\n", "\n"))
	}

	assert.Equal(t, string(expected5), string(gen))
}

func TestWrite32(t *testing.T) {
	doTestWrite(t, "386")
}

func TestWrite64(t *testing.T) {
	doTestWrite(t, "amd64")
}

func TestWriteArm32(t *testing.T) {
	doTestWrite(t, "arm")
}

func TestWriteArm64(t *testing.T) {
	doTestWrite(t, "arm64")
}

func doTestWrite(t *testing.T, arch string) {
	filename := "cmd"

	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := os.ReadFile(path)
	assert.NoError(t, err)

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Error("Could not parse "+filename+".json", err)
	}
	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	tmpdir, err := os.MkdirTemp("", "resource")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	file := filepath.Join(tmpdir, "resource.syso")

	err = vi.WriteSyso(file, arch)
	assert.NoError(t, err)

	_, err = os.ReadFile(file)
	assert.NoError(t, err)
}

func TestMalformedJSON(t *testing.T) {
	filename := "bad"

	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := os.ReadFile(path)
	assert.NoError(t, err)

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config and return false
	if err := vi.ParseJSON(jsonBytes); err == nil {
		t.Error("Application was supposed to return error, got nil")
	}
}

// TestSysoResourceAlignment guards against a regression of
// https://github.com/josephspurrier/goversioninfo/issues/39. Every resource in
// the .rsrc section must start on an 8 byte boundary. When it does not, the
// binutils linker used by mingw-w64 rejects the object file while merging
// resources with ".rsrc merge failure: corrupt .rsrc section". Odd sized
// resources such as an icon group (6+14*n bytes) are what push the following
// resources out of alignment, so the fixture below uses a multi image icon.
func TestSysoResourceAlignment(t *testing.T) {
	path, _ := filepath.Abs("./testdata/json/cmd.json")

	jsonBytes, err := os.ReadFile(path)
	assert.NoError(t, err)

	vi := &VersionInfo{}
	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Fatal("Could not parse cmd.json", err)
	}

	vi.IconPath = "testdata/resource/icon.ico"
	vi.ManifestPath = "testdata/resource/goversioninfo.exe.manifest"

	vi.Build()
	vi.Walk()

	for _, arch := range []string{"386", "amd64", "arm", "arm64"} {
		t.Run(arch, func(t *testing.T) {
			tmpdir, err := os.MkdirTemp("", "alignment")
			assert.NoError(t, err)
			defer os.RemoveAll(tmpdir)

			file := filepath.Join(tmpdir, "resource.syso")
			assert.NoError(t, vi.WriteSyso(file, arch))

			f, err := pe.Open(file)
			if err != nil {
				t.Fatal("Could not open the generated .syso", err)
			}
			defer f.Close()

			section := f.Section(".rsrc")
			if section == nil {
				t.Fatal("Generated .syso has no .rsrc section")
			}
			data, err := section.Data()
			assert.NoError(t, err)

			for _, leaf := range resourceLeaves(t, data, section.VirtualAddress) {
				if leaf%8 != 0 {
					t.Errorf("resource at %#x is not 8 byte aligned (offset mod 8 = %d)", leaf, leaf%8)
				}
			}

			// The offsets recorded in the headers must agree with what was
			// actually written, otherwise the padding was accounted for during
			// Freeze but left out of the file (or the other way around).
			if got, want := section.Offset+section.Size, section.PointerToRelocations; got != want {
				t.Errorf("raw data ends at %#x, relocations start at %#x", got, want)
			}
			relocationsEnd := section.PointerToRelocations + uint32(section.NumberOfRelocations)*10
			if got, want := relocationsEnd, f.FileHeader.PointerToSymbolTable; got != want {
				t.Errorf("relocations end at %#x, symbol table starts at %#x", got, want)
			}
		})
	}
}

// resourceLeaves walks the resource directory tree in the raw .rsrc data and
// returns the section relative offset of every resource's data.
func resourceLeaves(t *testing.T, data []byte, virtualAddress uint32) []uint32 {
	t.Helper()

	var offsets []uint32

	// read returns the little endian uint32 at off, failing the test when the
	// generated file is too short to hold it.
	read := func(off uint32) uint32 {
		if uint64(off)+4 > uint64(len(data)) {
			t.Fatalf("offset %#x is outside the %d byte .rsrc section", off, len(data))
		}
		return binary.LittleEndian.Uint32(data[off:])
	}

	// The tree is type -> ID -> language, so it is never deeper than 3 levels.
	var walk func(dirOffset uint32, depth int)
	walk = func(dirOffset uint32, depth int) {
		if depth > 3 {
			t.Fatalf("resource directory nested more than 3 levels deep at %#x", dirOffset)
		}

		// The entry count is the last uint32 of the 16 byte directory header:
		// the named entry count followed by the ID entry count.
		counts := read(dirOffset + 12)
		total := (counts & 0xffff) + (counts >> 16)

		for i := uint32(0); i < total; i++ {
			entry := dirOffset + 16 + i*8
			offsetToData := read(entry + 4)
			if offsetToData&0x80000000 != 0 {
				walk(offsetToData&0x7fffffff, depth+1)
				continue
			}
			// A leaf points at an IMAGE_RESOURCE_DATA_ENTRY, whose first field
			// is the RVA of the resource itself.
			offsets = append(offsets, read(offsetToData)-virtualAddress)
		}
	}
	walk(0, 1)

	if len(offsets) == 0 {
		t.Fatal("Found no resources in the .rsrc section")
	}

	return offsets
}

func TestIcon(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := os.ReadFile(path)
	assert.NoError(t, err)

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Error("Could not parse "+filename+".json", err)
	}

	vi.IconPath = "testdata/resource/icon.ico"

	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	tmpdir, err := os.MkdirTemp("", "resource")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	file := filepath.Join(tmpdir, "resource.syso")

	err = vi.WriteSyso(file, "386")
	assert.NoError(t, err)

	_, err = os.ReadFile(file)
	assert.NoError(t, err)
}

func TestBadIcon(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := os.ReadFile(path)
	assert.NoError(t, err)

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Error("Could not parse "+filename+".json", err)
	}

	vi.IconPath = "icon2.ico"

	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	tmpdir, err := os.MkdirTemp("", "resource")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	file := filepath.Join(tmpdir, "resource.syso")

	err = vi.WriteSyso(file, "386")
	if err == nil {
		t.Errorf("Error is missing because it should throw an error")
	}

	_, err = os.ReadFile(file)
	if err == nil {
		t.Error("File should not exist "+file, err)
	}
}

func TestApplicationIconDefault(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := os.ReadFile(path)
	assert.NoError(t, err)

	vi := &VersionInfo{}

	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Error("Could not parse "+filename+".json", err)
	}

	vi.IconPath = "testdata/resource/icon.ico"

	vi.Build()
	vi.Walk()

	tmpdir, err := os.MkdirTemp("", "resource")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	file := filepath.Join(tmpdir, "resource.syso")

	err = vi.WriteSyso(file, "386")
	assert.NoError(t, err)

	_, err = os.ReadFile(file)
	assert.NoError(t, err)
}

func TestApplicationIconExplicit(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := os.ReadFile(path)
	assert.NoError(t, err)

	vi := &VersionInfo{}

	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Error("Could not parse "+filename+".json", err)
	}

	vi.IconPath = "testdata/resource/icon.ico"
	vi.ApplicationIconPath = "testdata/resource/icon.ico"

	vi.Build()
	vi.Walk()

	tmpdir, err := os.MkdirTemp("", "resource")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	file := filepath.Join(tmpdir, "resource.syso")

	err = vi.WriteSyso(file, "386")
	assert.NoError(t, err)

	_, err = os.ReadFile(file)
	assert.NoError(t, err)
}

func TestBadApplicationIcon(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := os.ReadFile(path)
	assert.NoError(t, err)

	vi := &VersionInfo{}

	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Error("Could not parse "+filename+".json", err)
	}

	vi.ApplicationIconPath = "nonexistent.ico"

	vi.Build()
	vi.Walk()

	tmpdir, err := os.MkdirTemp("", "resource")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	file := filepath.Join(tmpdir, "resource.syso")

	err = vi.WriteSyso(file, "386")
	assert.Error(t, err)
}

func TestTimestamp(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := os.ReadFile(path)
	assert.NoError(t, err)

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Error("Could not parse "+filename+".json", err)
	}

	vi.Timestamp = true

	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	tmpdir, err := os.MkdirTemp("", "resource")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	file := filepath.Join(tmpdir, "resource.syso")

	err = vi.WriteSyso(file, "386")
	assert.NoError(t, err)

	_, err = os.ReadFile(file)
	assert.NoError(t, err)
}

func TestVersionString(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := os.ReadFile(path)
	assert.NoError(t, err)

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Error("Could not parse "+filename+".json", err)
	}
	if vi.FixedFileInfo.GetVersionString() != "6.3.9600.16384" {
		t.Errorf("Version String does not match: %v", vi.FixedFileInfo.GetVersionString())
	}
}

func TestWriteHex(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := os.ReadFile(path)
	assert.NoError(t, err)

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config
	if err := vi.ParseJSON(jsonBytes); err != nil {
		t.Error("Could not parse "+filename+".json", err)
	}
	// Fill the structures with config data
	vi.Build()

	// Write the data to a buffer
	vi.Walk()

	tmpdir, err := os.MkdirTemp("", "resource")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	file := filepath.Join(tmpdir, "resource.syso")

	err = vi.WriteHex(file)
	assert.NoError(t, err)

	_, err = os.ReadFile(file)
	assert.NoError(t, err)
}

func TestStr2Uint32(t *testing.T) {
	for _, tt := range []struct {
		in  string
		out uint32
	}{{"0", 0}, {"", 0}, {"FFEF", 65519}, {"\x00\x00", 0}} {
		log.SetOutput(io.Discard)
		got := str2Uint32(tt.in)
		if got != tt.out {
			t.Errorf("%q: awaited %d, got %d.", tt.in, tt.out, got)
		}
	}
}

var unmarshals = []struct {
	in      string
	needErr bool
}{
	{"", false}, {"A", true}, {"1", false}, {`"FfeF"`, false},
	{`"FfeF`, true}, {`"FXXX"`, true},
}

func TestLangID(t *testing.T) {
	var lng LangID
	for _, tt := range unmarshals {
		if err := lng.UnmarshalJSON([]byte(tt.in)); tt.needErr && err == nil {
			t.Errorf("%q: needed error, got nil.", tt.in)
		} else if !tt.needErr && err != nil {
			t.Errorf("%q: got error: %v", tt.in, err)
		}
	}
}

func TestCharsetID(t *testing.T) {
	var cs CharsetID
	for _, tt := range unmarshals {
		if err := cs.UnmarshalJSON([]byte(tt.in)); tt.needErr && err == nil {
			t.Errorf("%q: needed error, got nil.", tt.in)
		} else if !tt.needErr && err != nil {
			t.Errorf("%q: got error: %v", tt.in, err)
		}
	}
}

func TestWriteCoff(t *testing.T) {
	tempFh, err := os.CreateTemp("", "goversioninfo-test-")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	tempfn := tempFh.Name()
	tempFh.Close()
	defer os.Remove(tempfn)

	if err := writeCoff(nil, ""); err == nil {
		t.Errorf("needed error, got nil")
	}
	if err := writeCoff(nil, tempfn); err != nil {
		t.Errorf("got %v", err)
	}

	if err := writeCoffTo(badWriter{writeErr: io.EOF}, coff.NewRSRC()); err == nil {
		t.Errorf("needed write error, got nil")
	}
	if err := writeCoffTo(badWriter{closeErr: io.EOF}, nil); err == nil {
		t.Errorf("needed close error, got nil")
	}
}

func TestNewFileVersion(t *testing.T) {
	cases := []struct {
		in  string
		out FileVersion
		err string
	}{
		// Correct.
		{"1.2.3", FileVersion{1, 2, 3, 0}, ""},
		{"1.2.3.a", FileVersion{1, 2, 3, 0}, ""},
		{"1.2.3.4", FileVersion{1, 2, 3, 4}, ""},
		{"1.2.3.4-RC.1", FileVersion{1, 2, 3, 4}, ""},
		{"1.2.3.4 (final)", FileVersion{1, 2, 3, 4}, ""},
		{"6.3.9600.17284 (aaa.140822-1915)", FileVersion{6, 3, 9600, 17284}, ""},

		// Unexpected format.
		{"1.2", FileVersion{}, "version expected to start from x.y.z"},
		{"1.3.a", FileVersion{}, "version expected to start from x.y.z"},
		{"v1.2.3", FileVersion{}, "version expected to start from x.y.z"},

		// Any way to check Atoi errors except of overflow?
		{"1.1.1.9223372036854775808", FileVersion{}, "9223372036854775808"},
	}
	for i, c := range cases {
		got, err := NewFileVersion(c.in)
		if err == nil && c.err == "" && got != c.out {
			t.Errorf("%d) %q: expected %+v got %+v", i, c.in, c.out, got)
		} else if err == nil && c.err != "" {
			t.Errorf("%d) %q: expected error with susbtring %q got nil", i, c.in, c.err)
		} else if err != nil && c.err == "" {
			t.Errorf("%d) %q: unexpected error %s", i, c.in, err)
		} else if err != nil && c.err != "" && !strings.Contains(err.Error(), c.err) {
			t.Errorf("%d) %q: expected error with susbtring %q got %s", i, c.in, c.err, err)
		}
	}
}

func TestFillVersions_FixedToString(t *testing.T) {
	vi := &VersionInfo{}
	vi.FixedFileInfo.FileVersion = FileVersion{2, 0, 0, 0}
	vi.FixedFileInfo.ProductVersion = FileVersion{3, 1, 0, 0}
	vi.Build()
	assert.Equal(t, "2.0.0.0", vi.StringFileInfo.FileVersion)
	assert.Equal(t, "3.1.0.0", vi.StringFileInfo.ProductVersion)
}

func TestFillVersions_StringToFixed(t *testing.T) {
	vi := &VersionInfo{}
	vi.StringFileInfo.FileVersion = "2.0.0.0"
	vi.StringFileInfo.ProductVersion = "3.1.4.1"
	vi.Build()
	assert.Equal(t, FileVersion{2, 0, 0, 0}, vi.FixedFileInfo.FileVersion)
	assert.Equal(t, FileVersion{3, 1, 4, 1}, vi.FixedFileInfo.ProductVersion)
}

func TestFillVersions_BothPresent(t *testing.T) {
	vi := &VersionInfo{}
	vi.FixedFileInfo.FileVersion = FileVersion{1, 0, 0, 0}
	vi.FixedFileInfo.ProductVersion = FileVersion{1, 0, 0, 0}
	vi.StringFileInfo.FileVersion = "1.0.0.0 (custom build)"
	vi.StringFileInfo.ProductVersion = "1.0"
	vi.Build()
	assert.Equal(t, "1.0.0.0 (custom build)", vi.StringFileInfo.FileVersion)
	assert.Equal(t, "1.0", vi.StringFileInfo.ProductVersion)
	assert.Equal(t, FileVersion{1, 0, 0, 0}, vi.FixedFileInfo.FileVersion)
	assert.Equal(t, FileVersion{1, 0, 0, 0}, vi.FixedFileInfo.ProductVersion)
}

func TestFillVersions_InvalidStringVersion(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	vi := &VersionInfo{}
	vi.StringFileInfo.FileVersion = "not-a-version"
	vi.Build()
	assert.Equal(t, FileVersion{0, 0, 0, 0}, vi.FixedFileInfo.FileVersion)
	assert.Equal(t, "not-a-version", vi.StringFileInfo.FileVersion)
	assert.Contains(t, buf.String(), "could not be parsed")
}

func TestFillVersions_Conflict(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	vi := &VersionInfo{}
	vi.FixedFileInfo.FileVersion = FileVersion{2, 0, 0, 0}
	vi.StringFileInfo.FileVersion = "3.0.0.0"
	vi.Build()
	assert.Contains(t, buf.String(), "do not match")
}

func TestFillVersions_MatchingVersionsNoWarning(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	vi := &VersionInfo{}
	vi.FixedFileInfo.FileVersion = FileVersion{6, 3, 9600, 16384}
	vi.StringFileInfo.FileVersion = "6.3.9600.16384 (winblue_rtm.130821-1623)"
	vi.Build()
	assert.Empty(t, buf.String())
}

func TestPadStringEmoji(t *testing.T) {
	// U+1F600 (😀) requires a UTF-16 surrogate pair: 0xD83D 0xDE00
	got := padString("😀", 0)
	expected := []byte{0x3D, 0xD8, 0x00, 0xDE}
	assert.Equal(t, expected, got)
}

func TestPadStringMixed(t *testing.T) {
	// Mix of ASCII, non-ASCII BMP (©, U+00A9), and supplementary plane (😀, U+1F600)
	got := padString("a©😀", 0)
	expected := []byte{
		0x61, 0x00, // 'a'
		0xA9, 0x00, // '©'
		0x3D, 0xD8, 0x00, 0xDE, // '😀'
	}
	assert.Equal(t, expected, got)
}

func TestCopyrightSymbolCharset(t *testing.T) {
	copyright := "© 2024 Contoso Ltd."

	// padString always produces UTF-16LE regardless of CharsetID.
	encoded := padString(copyright, 0)
	// First two bytes should be © (U+00A9) encoded as little-endian UTF-16.
	assert.Equal(t, byte(0xA9), encoded[0], "low byte of © should be 0xA9")
	assert.Equal(t, byte(0x00), encoded[1], "high byte of © should be 0x00")

	// Build two resources with the same copyright string but different charsets.
	build := func(charset CharsetID) []byte {
		vi := &VersionInfo{}
		vi.StringFileInfo.LegalCopyright = copyright
		vi.VarFileInfo.Translation.LangID = LangID(0x0409)
		vi.VarFileInfo.Translation.CharsetID = charset
		vi.Build()
		vi.Walk()
		return vi.Buffer.Bytes()
	}

	correct := build(CsUnicode) // 04B0 — matches the UTF-16LE encoding
	wrong := build(Cs7ASCII)    // 0000 — declares 7-bit ASCII

	// The copyright bytes are identical in both buffers because padString
	// always writes UTF-16LE. The only difference is the Translation metadata
	// that tells Windows how to interpret those bytes.
	assert.Contains(t, string(correct), string(encoded),
		"correct charset buffer should contain UTF-16LE copyright bytes")
	assert.Contains(t, string(wrong), string(encoded),
		"wrong charset buffer should also contain the same UTF-16LE copyright bytes")

	// The buffers must differ — the Translation metadata is different.
	assert.NotEqual(t, correct, wrong,
		"buffers should differ because Translation charset metadata differs")

	// Verify the translation strings reflect the charset choice.
	correctTrans := Translation{LangID: 0x0409, CharsetID: CsUnicode}
	wrongTrans := Translation{LangID: 0x0409, CharsetID: Cs7ASCII}
	assert.Equal(t, "040904B0", correctTrans.getTranslationString(),
		"correct translation should declare Unicode (04B0)")
	assert.Equal(t, "04090000", wrongTrans.getTranslationString(),
		"wrong translation declares 7-bit ASCII (0000) — Windows will show '?' for non-ASCII characters")
}

type badWriter struct {
	writeErr, closeErr error
}

func (w badWriter) Write(p []byte) (int, error) {
	return len(p), w.writeErr
}
func (w badWriter) Close() error {
	return w.closeErr
}
