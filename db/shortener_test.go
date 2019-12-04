package db_test

import (
	"fmt"
	"testing"

	"github.com/zhanbei/static-server/db"
	"github.com/zhanbei/static-server/shortener"
)

// Package --> Models
// Package --> Handlers
func TestLoadShortenerRedirectionRecords(t *testing.T) {
	records := make(shortener.Records, 0)
	err := db.LoadShortenerRedirectionRecords(&records)
	fmt.Println(err, records)
}
