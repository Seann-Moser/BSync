package cmd

import (
	"github.com/Seann-Moser/BSync/internal/search"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var songsSearchCmd = &cobra.Command{
	Use:   "song-search",
	Short: "",
	Long:  "",
	Run:   search.Runner,
}

func songsSearchInit() {
	songsSearchCmd.Flags().StringP("search", "s", "", "")
	viper.BindPFlags(songsSearchCmd.Flags())
	RootCmd.AddCommand(songsSearchCmd)
}
