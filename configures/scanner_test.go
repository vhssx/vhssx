package configures_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/zhanbei/static-server/configures"
)

func TestScanSites(t *testing.T) {
	global, modular, sites, err := configures.ScanSites("../demo")
	fmt.Println("\nerr:", err)
	fmt.Println("global:", str(global))
	fmt.Println("modular:", str(modular))
	fmt.Println("sites:", str(sites))
}

func str(v interface{}) string {
	bts, _ := json.Marshal(v)
	return string(bts)
}
