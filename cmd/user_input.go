package cmd

import (
	"fmt"
	"github.com/Seann-Moser/BSync/internal/configuration"
	"github.com/Seann-Moser/BSync/internal/parser"
	"go.uber.org/zap"
	"log"
	"time"
)

func UserInput() {
	config, err := configuration.LoadConfigWithUserInput()
	if err != nil {
		log.Printf(err.Error())
		return
	}
	songData, err := parser.LoadSongEndpoints(config)
	if err != nil {
		config.Logger.Error("failed loading song config", zap.Error(err))
		time.Sleep(10 * time.Second)
		return
	}
	for _, s := range songData {
		err = s.Process(*config)
		if err != nil {
			config.Logger.Error("failed downloading songs", zap.Error(err))
			time.Sleep(10 * time.Second)
		}
	}

	config.Logger.Info("Press return to continue...")
	fmt.Scanln()

}
