package helper

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func deletePath(t *testing.T, path string) {
	err := os.RemoveAll(path)
	assert.NoError(t, err, "failed to delete temporary directory")
}

func TestMoveFile(t *testing.T) {
	type args struct {
		sourcePath string
		createPath string
		destPath   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"same-directory rename", args{"myfile", "myfile", "myfile2"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmp, err := ioutil.TempDir(os.TempDir(), "myprefix")
			assert.NoError(t, err, "failed to create temp directory")
			defer deletePath(t, tmp)

			assert.NoError(t, os.Chdir(tmp), "failed to change to testing directory")

			Execute("touch {{.}}", tt.args.createPath)
			assert.NoError(t, ioutil.WriteFile(tt.args.createPath, []byte("hello"), 0755), "failed to write test data to creation path")

			// Doing files not directories.
			//assert.NoError(t, os.MkdirAll(creationPath, 0755), "failed to create directory to move")

			want, err := ioutil.ReadFile(tt.args.createPath)
			assert.NoError(t, err, "createPath should exist befire MoveFile()")

			if err := MoveFile(tt.args.sourcePath, tt.args.destPath); (err != nil) != tt.wantErr {
				t.Errorf("MoveFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !DoesPathExist(tt.args.destPath) {
				t.Errorf("MoveFile() file at destPath %s not found", tt.args.destPath)
			}

			if !DoesPathExist(tt.args.createPath) {
				t.Errorf("MoveFile() copy did remove create path %s", tt.args.destPath)
			}

			want, err = ioutil.ReadFile(tt.args.createPath)
			assert.NoError(t, err, "error while reading form create path")

			got, err := ioutil.ReadFile(tt.args.destPath)
			assert.NoError(t, err, "failed to read from destination file path")

			assert.Equal(t, want, got, "source and destination files are not equal")
		})
	}
}
