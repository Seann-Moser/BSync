package search

import (
	"fmt"
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
	u := "https://bsaber.com/?s=" + strings.ReplaceAll(config.Search, " ", "+")
	songs, err := songParser.GetSongsWithPage(u, config.SongDownloadAmount, config.MinRatingPercent)
	if err != nil {
		config.Logger.Fatal("failed getting songs from page:"+config.BSaberURL, zap.Error(err))
	}
	config.Logger.Info(fmt.Sprintf("found %d songs for this search", len(songs)))
	songParser.DownloadSongList(songs, config.Workers, config.BeatSaberPath)
}
