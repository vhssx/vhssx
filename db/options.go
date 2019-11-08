package db

const (
	ColRequests = "requests"
)

type MongoDbOptions struct {
	Uri string `json:"uri"`

	DbName string `json:"database"`
	// For all normal records/requests.
	// logging.normal.requests
	// logging.normal.devices
	// logging.normal.clicks
	// logging.normal.views
	CollectionPrefix string `json:"collection"`
}

func (m *MongoDbOptions) GetColName(colName string) string {
	return m.CollectionPrefix + "." + colName
}
