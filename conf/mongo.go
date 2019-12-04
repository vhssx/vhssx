package conf

const (
	// Logger for first-level/unvalidated requests
	PrefixColRequests = "requests"
	// Logger for first-level/general requests
	ColCrawlerRequests = PrefixColRequests + ".crawlers"
	// Col for the general(first-level) requests
	ColGeneralRequests = PrefixColRequests + ".general"
	// Col for the first-level requests.
	ColLandingRequests = PrefixColRequests + ".landing"
	// Col for the second-level/validated requests.
	ColValidatingRequests = PrefixColRequests + ".validating"
	// Col for the third-level/post-validated requests.
	ColValidatedRequests = PrefixColRequests + ".validated"
	// Unknown requests.
	ColUnknownRequests = PrefixColRequests + ".unknown"

	// URL Shortener / Redirection Management
	PrefixColShortener = "shortener"
	// The records of all redirections.
	ColShortenerRedirections = PrefixColShortener + ".redirections"

	// Store the configures later.
	PrefixColConfigures = "configures"
	// The configures of domains.
	// FIX-ME Support configures for domains from database( vs configuration file).
	ColConfigureDomains = PrefixColConfigures + ".domains"
)

type MongoDbOptions struct {
	Enabled bool `json:"enabled"`

	Uri string `json:"uri"`

	DbName string `json:"db"`
	// For all normal records/requests.
	// logging.normal.requests
	// logging.normal.devices
	// logging.normal.clicks
	// logging.normal.views
	CollectionPrefix string `json:"col"`
}

func (m *MongoDbOptions) GetColName(colName string) string {
	return m.CollectionPrefix + "." + colName
}

func (m *MongoDbOptions) IsValid() bool {
	if doPass(m.Enabled) {
		return true
	}
	return exist(m.Uri) && exist(m.DbName) && exist(m.CollectionPrefix)
}
