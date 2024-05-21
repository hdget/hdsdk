package sqlboiler_sqlite3

import (
	"github.com/hdget/hdsdk/v2/errdef"
	"github.com/hdget/hdsdk/v2/intf"
	"github.com/pkg/errors"
)

type sqliteProviderConfig struct {
	DbName string `mapstructure:"dbname"`
}

const (
	configSection = "sdk.sqlite"
)

func newConfig(configProvider intf.ConfigProvider) (*sqliteProviderConfig, error) {
	if configProvider == nil {
		return nil, errdef.ErrInvalidConfig
	}

	var c *sqliteProviderConfig
	err := configProvider.Unmarshal(&c, configSection)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, errdef.ErrEmptyConfig
	}

	err = c.validate()
	if err != nil {
		return nil, errors.Wrap(err, "validate mysql provider config")
	}

	return c, nil
}

func (c *sqliteProviderConfig) validate() error {
	if c == nil || c.DbName == "" {
		return errdef.ErrInvalidConfig
	}
	return nil
}
