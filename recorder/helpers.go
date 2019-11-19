package recorder

import (
	"encoding/json"
	"fmt"
)

// Print out the record, once there are failures when serializing recording.
func PrintFailedRecordText(record string) {
	fmt.Println("-+/>", record)
}

// Print out the record, once there are failures when serializing recording.
func PrintFailedRecord(record interface{}) {
	bts, _ := json.Marshal(record)
	PrintFailedRecordText(string(bts))
}
