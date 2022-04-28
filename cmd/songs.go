package cmd

import (
	"github.com/Seann-Moser/BSync/internal/songs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var songsCmd = &cobra.Command{
	Use:   "song",
	Short: "",
	Long:  "",
	Run:   songs.Runner,
}

func songsInit() {
	songsCmd.Flags().StringP("beat-sync-url", "b", "https://bsaber.com/songs/top/?time=7-days", "")
	viper.BindPFlags(songsCmd.Flags())
	RootCmd.AddCommand(songsCmd)
}
