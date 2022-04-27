package configuration

import (
	"time"

	"github.com/Netflix/go-env"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Config struct {
	Port string `env:"PORT,default=8080"`

	DBHost     string `env:"DB_HOST,default=mysql"`
	DBUser     string `env:"DB_USER,default=root"`
	DBPassword string `env:"DB_PASSWORD,default=root"`
	DBPort     string `env:"DB_PORT,default=3308"`

	TIMEOUT time.Duration `env:"TIMEOUT,default=30s"`
	Version string        `env:"VERSION"`
	Logger  *zap.Logger

	LoggingLevel string `env:"LOGGING_LEVEL,default=DEBUG"`
	LoggingProd  bool   `env:"LOGGING_PROD,default=true"`
	Extras       env.EnvSet

	DB        *sqlx.DB
	RateLimit int `env:"RATE_LIMIT,default=120"`
}

func LoadConfig() (*Config, error) {
	var NewConfig Config
	es, err := env.UnmarshalFromEnviron(&NewConfig)
	if err != nil {
		return nil, err
	}
	NewConfig.Extras = es
	NewConfig.Logger, err = ConfigureLogger(&NewConfig)
	if err != nil {
		return nil, err
	}
	err = connectToDB(&NewConfig)
	if err != nil {
		return nil, err
	}
	return &NewConfig, nil
}
