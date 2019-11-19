package conf

const (
	ColRequests = "requests"
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
