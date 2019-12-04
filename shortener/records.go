package shortener

import (
	"github.com/zhanbei/dxb"
)

type Records = []*Record

// Case-(in)sensitive Capacity:
// (26*2)^6 = 26^6 * 2^6
// (26*2 + 10)^6 >> (26 + 10)^6
type Record struct {
	Id dxb.ObjectId `json:"_id" bson:"_id"`
	// Active or not.
	Enabled bool `json:"enabled"`
	// The domain of the short URL.
	// FIX-ME Extract the #Domain out and use the #DomainId as the foreign ID.
	Domain string `json:"domain"`
	// The prefix of the key.
	// The default is "/", can be set to "/s/" to improve the performance.
	Prefix string `json:"prefix"`
	// The key should be unique, or the later record will override the previous one.
	Key string `json:"key"`
	// Exactly match > Case Insensitive Match > 404 Not Found
	// Generally, the feature should be on, when the link is possibly typed by the end users.
	// The feature should be off, to be secure and have large capacity.
	// The default false means sensitive.
	CaseInsensitive bool `json:"insensitive" bson:"insensitive"`

	Title string `json:"title"`
	// The destination (full )URL with domains and HTTP(S) protocol.
	Target string `json:"url"`

	Tags []string `json:"tags"`
	// A valida HTTP code for redirection, or an invalid one to use the default code.
	// Record Code > User Default Code > Platform Default Code
	Code int `json:"code"`
}
