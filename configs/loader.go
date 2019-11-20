package configs

import (
	"errors"

	"github.com/micro/go-micro/config"
	"github.com/zhanbei/static-server/conf"
	"github.com/zhanbei/static-server/utils"
)

// FIX-ME Think about dynamic hot loader.
func LoadServerConfigures(mGivenConfigFile string, rawServerOptions *conf.ServerOptions, rawAddress, rawRootDir string) (*conf.Configure, error) {
	if mGivenConfigFile == "" {
		return nil, errors.New("expected configuration file")
	}

	err := config.LoadFile(mGivenConfigFile)
	if err != nil {
		return nil, err
	}

	mConfig := new(conf.Configure)

	err = config.Scan(mConfig)
	if err != nil {
		return mConfig, err
	}

	if utils.NotEmpty(rawAddress) {
		mConfig.Address = rawAddress
	}
	if utils.NotEmpty(rawRootDir) {
		mConfig.RootDir = rawRootDir
	}
	if mConfig.Server == nil {
		// FIX-ME Prefer the configures in the configuration file over cli arguments.
		mConfig.Server = rawServerOptions
	}

	if !mConfig.IsValid() {
		return mConfig, errors.New("invalid configures, see help for more info")
	}

	return mConfig, nil
}
