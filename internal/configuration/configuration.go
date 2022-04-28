package configuration

import (
	"encoding/json"
	"fmt"
	"github.com/Netflix/go-env"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strings"
)

const (
	BaseConfigName = "user_config.json"
)

type Config struct {
	Logger *zap.Logger `json:"-"`

	LoggingLevel string `mapstructure:"logging-level" env:"LOGGING_LEVEL,default=info"`
	LoggingProd  bool   `mapstructure:"logging-prod" env:"LOGGING_PROD,default=true"`
	BSaberURL    string `mapstructure:"beat-sync-url"`
	Workers      int    `mapstructure:"workers" env:"WORKERS,default=4"`

	Search             string `json:"search" mapstructure:"search"`
	UserName           string `json:"user_name" mapstructure:"user_name" env:"USER_NAME,default="`
	SongDownloadAmount int    `mapstructure:"song-download-amount" env:"SONG_AMOUNT,default=20"`
	BeatSaberPath      string `mapstructure:"beat-saber-path" env:"BEAT_SABER_PATH,default=C:/Program Files (x86)/Steam/steamapps/common/Beat Saber/Beat Saber_Data/CustomLevels/"`
	DownloadDelay      int    `mapstructure:"download-delay" env:"DOWNLOAD_DELAY,default=5"`

	MinRatingPercent float32 `mapstructure:"min-rating-percent" env:"MIN_RATING_PERCENT,default=0.5"`
	FirstLoad        bool    `env:"first_load,default=true"`
	SongConfigName   string  `env:"SONG_CONFIG_NAME.json,default=songs_config.json"`
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

func LoadConfigWithUserInput() (*Config, error) {
	var conf Config
	if e, err := exists(BaseConfigName); err == nil && e {
		file, err := ioutil.ReadFile(BaseConfigName)
		if err == nil {
			err = json.Unmarshal([]byte(file), &conf)
			if err != nil {
				return nil, err
			}
			conf.Logger, err = ConfigureLogger(&conf)
			if err != nil {
				return nil, err
			}
			return &conf, err
		}
	}

	_, err := env.UnmarshalFromEnviron(&conf)
	if err != nil {
		return nil, err
	}
	conf.Logger, err = ConfigureLogger(&conf)
	if err != nil {
		return nil, err
	}
	if conf.FirstLoad {
		fmt.Print("Enter BSaber(https://bsaber.com/) User Name: ")
		_, err = fmt.Scanln(&conf.UserName)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Enter Beat Saber Path(Default:%s) User Name: ", conf.BeatSaberPath)
		path := ""
		_, err = fmt.Scanln(&path)
		if err != nil && !strings.Contains(err.Error(), "unexpected newline") {
			return nil, err
		}
		if path != "" {
			if exists, err := exists(path); err != nil || !exists {
				conf.Logger.Error(fmt.Sprintf("path: %s does not exist", path), zap.Error(err))
				return nil, err
			}
			conf.BeatSaberPath = path
		}
	}
	conf.FirstLoad = false
	file, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(BaseConfigName, file, 0644)
	return &conf, err
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
