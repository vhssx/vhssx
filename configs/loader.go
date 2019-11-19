package configs

import (
	"errors"

	"github.com/micro/go-micro/config"
	"github.com/zhanbei/static-server/conf"
)

func LoadServerConfigures(mGivenConfigFile string) (*conf.Configure, error) {
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
	if !mConfig.IsValid() {
		return mConfig, errors.New("invalid configures, see help for more info")
	}

	err = mConfig.ValidateFile()
	if err != nil {
		return mConfig, errors.New("logical error with valid configures")
	}

	return mConfig, nil
}
