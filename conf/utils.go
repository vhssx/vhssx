package conf

import (
	"os"
)

// Get the file writer for file loggers.
func GetFileWriter(target string) (*os.File, error) {
	// @see https://stackoverflow.com/questions/7151261/append-to-a-file-in-go
	return os.OpenFile(target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}
