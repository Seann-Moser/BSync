package songs

import (
	"github.com/Seann-Moser/BSync/internal/configuration"
	"github.com/Seann-Moser/BSync/internal/parser"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
)

func Runner(cmd *cobra.Command, args []string) {
	config, err := configuration.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	songParser := parser.NewSongParser(config.Logger)
	err = songParser.DownloadSongs(config.BSaberURL, config)
	if err != nil {
		config.Logger.Fatal("failed getting songs from page:"+config.BSaberURL, zap.Error(err))
	}

}
