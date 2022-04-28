package parser

import (
	"go.uber.org/zap"
	"testing"
)

func TestSongParser_GetSongs(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	//https://bsaber.com/songs/curated/?recommended=true
	//https://bsaber.com/songs/top/?time=30-days
	//https://bsaber.com/songs/top/?time=7-days
	parser := NewSongParser(logger)
	_, err = parser.GetSongsWithPage("https://bsaber.com/songs/curated/?recommended=true", 25, 0.5, "")
	if err != nil {
		t.Fatal(err)
	}

	//parser.DownloadSongList(songs, 4, "C:/Program Files (x86)/Steam/steamapps/common/Beat Saber/Beat Saber_Data/CustomLevels/")
}
