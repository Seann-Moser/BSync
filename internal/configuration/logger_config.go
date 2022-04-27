package configuration

import "go.uber.org/zap"

func ConfigureLogger(conf *Config) (*zap.Logger, error) {
	var loggerConfig zap.Config
	if conf.LoggingProd {
		loggerConfig = zap.NewProductionConfig()
	} else {
		loggerConfig = zap.NewDevelopmentConfig()
	}
	err := loggerConfig.Level.UnmarshalText([]byte(conf.LoggingLevel))
	if err != nil {
		return nil, err
	}

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
