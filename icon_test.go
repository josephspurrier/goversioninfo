package goversioninfo

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestIconReleaseFileHandle(t *testing.T) {
	icoPath := "testdata/resource/icon.ico"
	icoPath2 := "testdata/resource/icon2.ico"
	tmpdir, err := ioutil.TempDir("", "resource")
	if err != nil {
		t.Error("Could not create temp dir", err)
	}
	defer os.RemoveAll(tmpdir)
	outPath := filepath.Join(tmpdir, "resource.syso")
	vi := &VersionInfo{}
	vi.IconPath = icoPath

	vi.Build()
	vi.Walk()
	err = vi.WriteSyso(outPath, runtime.GOARCH)
	if err != nil {
		t.Errorf("Unexpected error writing resource: %v", err)
	}

	err = os.Rename(icoPath, icoPath2)
	if err != nil {
		t.Errorf("Error renaming icon: %v", err)
	}

	err = os.Rename(icoPath2, icoPath)
	if err != nil {
		t.Errorf("Error restoring icon: %v", err)
	}
}
