package configuration

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Logger *zap.Logger

	LoggingLevel string `mapstructure:"logging-level"`
	LoggingProd  bool   `mapstructure:"logging-prod"`
	BSaberURL    string `mapstructure:"beat-sync-url"`
	Workers      int    `mapstructure:"workers"`

	Search string `json:"search"`

	SongDownloadAmount int    `mapstructure:"song-download-amount"`
	BeatSaberPath      string `mapstructure:"beat-saber-path"`
	DownloadDelay      int    `mapstructure:"download-delay"`

	MinRatingPercent float32 `mapstructure:"min-rating-percent"`
}

func LoadConfig() (*Config, error) {
	var conf Config
	err := viper.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}
	conf.Logger, err = ConfigureLogger(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
