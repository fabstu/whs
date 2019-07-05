package helper

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

func CurrentPackage() string {
	gomod, modPath, err := ModPathFromWD()
	if err != nil {
		panic(fmt.Sprintf("failed to find gomod to determine current package: %v", err))
	}

	currentDir := CurrentPackagePath()

	root := filepath.Dir(gomod)

	// Add module path in front, then add current file with root removed.
	return modPath + "/" + strings.TrimPrefix(currentDir, root+"/")
}

func CurrentPackagePath() string {
	_, currentFile, _, _ := runtime.Caller(1)
	return filepath.Dir(currentFile)
}