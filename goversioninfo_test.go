package goversioninfo

import (
	"bytes"
	"io"
	"io/ioutil"
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

	jsonBytes, err := ioutil.ReadFile(path)
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
	ioutil.WriteFile(path3, vi.Buffer.Bytes(), 0655)*/

	expected, err := ioutil.ReadFile(path2)
	assert.NoError(t, err)

	if !bytes.Equal(vi.Buffer.Bytes(), expected) {
		t.Error("Data does not match " + filename + ".hex")
	}

	// Test the Go file generation.
	tmpdir, err := ioutil.TempDir("", "generate_go")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	path4 := filepath.Join(tmpdir, filename+".go")
	err = vi.WriteGo(path4, "")
	assert.NoError(t, err)

	gen, err := ioutil.ReadFile(path4)
	assert.NoError(t, err)

	path5, _ := filepath.Abs("./testdata/gofile/" + filename + ".go")
	expected5, err := ioutil.ReadFile(path5)
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

func doTestWrite(t *testing.T, arch string) {
	filename := "cmd"

	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
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

	tmpdir, err := ioutil.TempDir("", "resource")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	file := filepath.Join(tmpdir, "resource.syso")

	err = vi.WriteSyso(file, arch)
	assert.NoError(t, err)

	_, err = ioutil.ReadFile(file)
	assert.NoError(t, err)
}

func TestMalformedJSON(t *testing.T) {
	filename := "bad"

	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
	assert.NoError(t, err)

	// Create a new container
	vi := &VersionInfo{}

	// Parse the config and return false
	if err := vi.ParseJSON(jsonBytes); err == nil {
		t.Error("Application was supposed to return error, got nil")
	}
}

func TestIcon(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
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

	tmpdir, err := ioutil.TempDir("", "resource")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	file := filepath.Join(tmpdir, "resource.syso")

	err = vi.WriteSyso(file, "386")
	assert.NoError(t, err)

	_, err = ioutil.ReadFile(file)
	assert.NoError(t, err)
}

func TestBadIcon(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
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

	tmpdir, err := ioutil.TempDir("", "resource")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	file := filepath.Join(tmpdir, "resource.syso")

	err = vi.WriteSyso(file, "386")
	if err == nil {
		t.Errorf("Error is missing because it should throw an error")
	}

	_, err = ioutil.ReadFile(file)
	if err == nil {
		t.Error("File should not exist "+file, err)
	}
}

func TestTimestamp(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
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

	tmpdir, err := ioutil.TempDir("", "resource")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	file := filepath.Join(tmpdir, "resource.syso")

	err = vi.WriteSyso(file, "386")
	assert.NoError(t, err)

	_, err = ioutil.ReadFile(file)
	assert.NoError(t, err)
}

func TestVersionString(t *testing.T) {
	filename := "cmd"

	path, _ := filepath.Abs("./testdata/json/" + filename + ".json")

	jsonBytes, err := ioutil.ReadFile(path)
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

	jsonBytes, err := ioutil.ReadFile(path)
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

	tmpdir, err := ioutil.TempDir("", "resource")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpdir)
	file := filepath.Join(tmpdir, "resource.syso")

	err = vi.WriteHex(file)
	assert.NoError(t, err)

	_, err = ioutil.ReadFile(file)
	assert.NoError(t, err)
}

func testdatatr2Uint32(t *testing.T) {
	for _, tt := range []struct {
		in  string
		out uint32
	}{{"0", 0}, {"", 0}, {"FFEF", 65519}, {"\x00\x00", 0}} {
		log.SetOutput(ioutil.Discard)
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
	tempFh, err := ioutil.TempFile("", "goversioninfo-test-")
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

type badWriter struct {
	writeErr, closeErr error
}

func (w badWriter) Write(p []byte) (int, error) {
	return len(p), w.writeErr
}
func (w badWriter) Close() error {
	return w.closeErr
}
