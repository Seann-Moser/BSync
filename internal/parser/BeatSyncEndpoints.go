package parser

import (
	"encoding/json"
	"fmt"
	"github.com/Seann-Moser/BSync/internal/configuration"
	"io/ioutil"
	"os"
)

type SongEndpoint struct {
	Url              string   `json:"url"`
	Amount           int      `json:"amount"`
	MinRating        float32  `json:"min_rating"`
	DifficultyLevels []string `json:"difficulty_levels"`
}

func (s *SongEndpoint) Process(conf configuration.Config) error {
	songParser := NewSongParser(conf.Logger)
	conf.SongDownloadAmount = s.Amount
	conf.MinRatingPercent = s.MinRating
	return songParser.DownloadSongs(s.Url, &conf)
}

func LoadSongEndpoints(conf *configuration.Config) ([]*SongEndpoint, error) {
	var output []*SongEndpoint
	if e, err := exists(conf.SongConfigName); err == nil && e {
		file, err := ioutil.ReadFile(conf.SongConfigName)
		if err == nil {
			err = json.Unmarshal(file, &output)
			return output, err
		}
	}
	output = setDefault(conf)
	file, err := json.MarshalIndent(output, "", "    ")
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(conf.SongConfigName, file, 0644)
	return output, err
}
func setDefault(conf *configuration.Config) []*SongEndpoint {
	return []*SongEndpoint{
		{
			Url:              fmt.Sprintf("https://bsaber.com/songs/new/?bookmarked_by=%s", conf.UserName),
			Amount:           -1, // Downloads all bookmarked songs
			MinRating:        0.0,
			DifficultyLevels: nil,
		},
		{
			Url:              "https://bsaber.com/songs/top/?time=30-days",
			Amount:           40,
			MinRating:        0.5,
			DifficultyLevels: nil,
		},
		{
			Url:              "https://bsaber.com/songs/top/?time=7-days",
			Amount:           20,
			MinRating:        0.5,
			DifficultyLevels: nil,
		},
		{
			Url:              "https://bsaber.com/songs/curated/?recommended=true",
			Amount:           20,
			MinRating:        0.5,
			DifficultyLevels: nil,
		},
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
