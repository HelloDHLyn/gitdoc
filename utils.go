package gitdoc

import (
	"os"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func stringOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
