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
	songsCmd.Flags().StringP("beat-sync-url", "u", "https://bsaber.com/songs/top/?time=7-days", "")
	songsCmd.Flags().IntP("workers", "w", 4, "")
	songsCmd.Flags().IntP("download-delay", "d", 5, "")
	songsCmd.Flags().IntP("song-download-amount", "a", 20, "")
	viper.BindPFlags(songsCmd.Flags())
	RootCmd.AddCommand(songsCmd)
}
