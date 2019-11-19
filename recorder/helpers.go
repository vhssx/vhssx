package recorder

import (
	"encoding/json"
	"fmt"
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

// Print out the record, once there are failures when serializing recording.
func PrintFailedRecordText(record string) {
	fmt.Println("-+/>", record)
}

// Print out the record, once there are failures when serializing recording.
func PrintFailedRecord(record interface{}) {
	bts, _ := json.Marshal(record)
	PrintFailedRecordText(string(bts))
}
