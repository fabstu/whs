package helper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// Finds the go mod file in the current or any directory further up.
func FindGoModFromWD() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get wd so did not find go mod: %v", err)
	}

	return FindGoModFrom(wd)
}

func FindGoModFrom(from string) (string, error) {
	possibleGomod := filepath.Join(from, "go.mod")
	if DoesPathExist(possibleGomod) {
		return possibleGomod, nil
	}

	nextToTry := filepath.Dir(from)
	if nextToTry == from {
		return "", fmt.Errorf("no go.mod in path")
	}

	return FindGoModFrom(nextToTry)
}

var (
	slashSlash = []byte("//")
	moduleStr  = []byte("module")
)

// Source: https://github.com/golang/go/blob/master/src/cmd/go/internal/modfile/read.go#L837

// ModulePath returns the module path from the gomod file text.
// If it cannot find a module path, it returns an empty string.
// It is tolerant of unrelated problems in the go.mod file.
func ModulePath(mod []byte) string {
	for len(mod) > 0 {
		line := mod
		mod = nil
		if i := bytes.IndexByte(line, '\n'); i >= 0 {
			line, mod = line[:i], line[i+1:]
		}
		if i := bytes.Index(line, slashSlash); i >= 0 {
			line = line[:i]
		}
		line = bytes.TrimSpace(line)
		if !bytes.HasPrefix(line, moduleStr) {
			continue
		}
		line = line[len(moduleStr):]
		n := len(line)
		line = bytes.TrimSpace(line)
		if len(line) == n || len(line) == 0 {
			continue
		}

		if line[0] == '"' || line[0] == '`' {
			p, err := strconv.Unquote(string(line))
			if err != nil {
				return "" // malformed quoted string or multiline module path
			}
			return p
		}

		return string(line)
	}
	return "" // missing module path
}
// Determines gomod path and modpath starting from the given path.
func ModPathFrom(from string) (string, string, error) {
	gomodPath, err := FindGoModFrom(from)
	if err != nil {
		return "", "", fmt.Errorf("failed find gomod file so release failed: %v", err)
	}

	moduleContent, err := ioutil.ReadFile(gomodPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read gomod file: %v", err)
	}
	return gomodPath, ModulePath([]byte(moduleContent)), nil
}

func ModPathFromWD() (string, string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", "", fmt.Errorf("failed to get wd so did not find mod path: %v", err)
	}

	return ModPathFrom(wd)
}