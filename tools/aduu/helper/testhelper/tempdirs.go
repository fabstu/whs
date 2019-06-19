package testhelper

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func MakeTempDir(t *testing.T, name string) string {
	tempDir, err := ioutil.TempDir(os.TempDir(), strings.ReplaceAll(name, " ", "-"))
	if err != nil {
		t.Fatal(err)
	}
	return tempDir
}

func DeleteTempDir(t *testing.T, path string) {
	if err := os.RemoveAll(path); err != nil {
		t.Fatal("failed to remove temp dir:", err)
	}
}
