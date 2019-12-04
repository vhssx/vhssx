package shortener

import (
	"fmt"
	"strings"

	"github.com/zhanbei/static-server/db"
)

// Get configures by domain.
type CacheDomainChecker map[string]bool
type CacheDomainMapper map[string]*Domain

func DoLoadRedirectionRecords() (CacheDomainChecker, CacheDomainMapper) {
	checker := make(CacheDomainChecker, 0)
	mapper := make(CacheDomainMapper, 0)

	records := LoadRedirectionRecords()
	for _, record := range records {
		if !checker[record.Domain] {
			checker[record.Domain] = true
			mapper[record.Domain] = &Domain{nil, make(RouteMapper, 0), true}
		}
		//routes := mapper[record.Domain].Routes
		mapper[record.Domain].Routes[record.Key] = record
		recordKey := strings.ToLower(record.Key)
		if record.CaseInsensitive && record.Key != recordKey {
			mapper[record.Domain].Routes[recordKey] = record
		}
	}

	return checker, mapper
}

func LoadRedirectionRecords() Records {
	records := make(Records, 0)
	err := db.LoadShortenerRedirectionRecords(&records)
	if err != nil {
		fmt.Println("Querying the records of redirection failed!")
	}
	fmt.Println("[[URL SHORTENER]] Found", len(records), "records of redirection!")
	return records
}
