package conf

import (
	"github.com/zhanbei/static-server/helpers/terminator"
	"github.com/zhanbei/static-server/utils"
)

type Configure struct {
	// The www-root-dir path.
	RootDir string `json:"rootDir"`
	// The address or port for the server.
	Address string `json:"address"`

	Server *ServerOptions `json:"server"`

	Loggers *OptionLoggers `json:"loggers"`

	MongoDbOptions *MongoDbOptions `json:"mongo"`

	GorillaOptions *OptionLoggerGorilla `json:"mongo"`
}

var has = utils.NotEmpty

func (m *Configure) IsValid() bool {
	if !has(m.RootDir) {
		terminator.ExitWithConfigError(nil, "Please specify an address( or at least a port) in your configuration file!")
		return false
	}
	m.Address, _ = ValidateArgAddressOrExit(m.Address)
	m.RootDir = ValidateArgRootDirOrExit(m.RootDir)

	return m.Server.IsValid()
}

func (m *Configure) ValidateFile() error {

	return nil
}
