package recorder

import (
	"io"
	"os"
)

func twoWriters(stdout bool, file *os.File) io.Writer {
	if !stdout {
		if file == nil {
			return nil
		} else {
			return file
		}
	} else {
		if file == nil {
			return os.Stdout
		} else {
			return io.MultiWriter(os.Stdout, file)
		}
	}
}
