package search

import (
	"github.com/Seann-Moser/BSync/internal/configuration"
	"github.com/Seann-Moser/BSync/internal/parser"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
	"strings"
)

func Runner(cmd *cobra.Command, args []string) {
	config, err := configuration.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	songParser := parser.NewSongParser(config.Logger)
	searchURL := "https://bsaber.com/?s=" + strings.ReplaceAll(config.Search, " ", "+")
	if config.Search == "" && config.UserName != "" {
		searchURL = "https://bsaber.com/songs/new/?bookmarked_by=" + config.UserName
	}
	err = songParser.DownloadSongs(searchURL, config)
	if err != nil {
		config.Logger.Fatal("failed getting songs from page:"+config.BSaberURL, zap.Error(err))
	}
}
