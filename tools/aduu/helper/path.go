package helper

import "os"

func DoesPathExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
