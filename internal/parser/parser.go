package parser

import (
	"encoding/json"
	"fmt"
	"github.com/Seann-Moser/BSync/internal/configuration"
	"github.com/Seann-Moser/WebParser/website"
	"go.uber.org/zap"
	"sync"
)

type SongParser struct {
	Data         []*website.Search `json:"id"`
	Difficulties []*website.Search `json:"difficulties"`
	Pages        []*website.Search `json:"pages"`
	Logger       *zap.Logger
	Process      *ParserProcessor
}

func NewSongParser(logger *zap.Logger) *SongParser {
	return &SongParser{
		Process: NewParserProcessor(logger),
		Data: []*website.Search{
			{
				Type:  website.TypeTag,
				Tag:   "article",
				Order: 0,
			},
			{
				Type:     website.TypeAttribute,
				Tag:      "class",
				TagValue: "post.*",
				Order:    0,
				Flatten:  true,
			},
			{
				Type:            website.TypeAttribute,
				Tag:             "data-original",
				InternalTagName: "thumbnail",
				Order:           0,
				OnlyRemap:       true,
			},
			{
				Type:            website.TypeAttribute,
				Tag:             "link",
				InternalTagName: "song_url",
				Order:           0,
				OnlyRemap:       true,
			},
			{
				Type:            website.TypeAttribute,
				Tag:             "href",
				InternalTagName: "raw_difficulties",
				Order:           0,
				OnlyRemap:       true,
			},
			{
				Type:            website.TypeAttribute,
				Tag:             "text",
				InternalTagName: "raw_text",
				Order:           0,
				OnlyRemap:       true,
			},
		},
		Difficulties: nil,
		Pages: []*website.Search{
			{
				Type:  website.TypeTag,
				Tag:   "div",
				Order: 0,
			},
			{
				Type:        website.TypeAttribute,
				Tag:         "class",
				TagValue:    "navigation pagination",
				Order:       0,
				Flatten:     false,
				ForwardData: true,
				SkipRemap:   true,
			},
			{
				Type:  website.TypeTag,
				Tag:   "a",
				Order: 1,
			},
		},
		Logger: logger,
	}
}

func (s *SongParser) DownloadSongs(u string, conf *configuration.Config) error {
	s.Logger.Info("grabbing songs from " + u)
	songs, err := s.GetSongsWithPage(u, conf.SongDownloadAmount, conf.MinRatingPercent, conf.BeatSaberPath)
	if err != nil {
		return err
	}
	if len(songs) > 0 {
		s.DownloadSongList(songs, conf.Workers, conf.BeatSaberPath)
		s.Logger.Info(fmt.Sprintf("finished downloading %d new songs", len(songs)))
	} else {
		s.Logger.Info("no songs found")
	}

	return nil
}

func (s *SongParser) GetSongsWithPage(u string, amount int, minRating float32, path string) ([]*Song, error) {
	var songList []*Song
	currentURL := u
	visitedMap := map[string]bool{}
	for {
		visitedMap[currentURL] = true
		sl, err := s.GetSongs(currentURL, minRating, path)
		if err != nil {
			return nil, err
		}

		songList = append(songList, sl...)
		if len(songList) > amount {
			return songList[:amount], nil
		}

		pageData, err := s.Process.GetData(currentURL, s.Pages)
		if err != nil {
			return nil, err
		}
		type Page struct {
			Link string `json:"link"`
			Text string `json:"text"`
		}
		p := []*Page{}
		err = json.Unmarshal(pageData, &p)
		if err != nil {
			return nil, err
		}
		if len(p) > 0 {
			foundLinks := false
			for _, i := range p {
				if _, found := visitedMap[i.Link]; !found && len(i.Text) > 0 && i.Text != "1" {
					currentURL = i.Link
					foundLinks = true
					break
				}
			}
			if !foundLinks {
				return songList, err
			}
		} else {
			return songList, err
		}
	}

}
func (s *SongParser) GetSongs(u string, minRating float32, path string) ([]*Song, error) {
	b, err := s.Process.GetData(u, s.Data)
	if err != nil {
		return nil, err
	}
	var songList []*Song
	err = json.Unmarshal(b, &songList)
	if err != nil {
		return nil, err
	}
	output := []*Song{}
	for _, s := range songList {
		s.Process()
		if !s.AlreadyDownloaded(path) && s.RatingPercent >= minRating {
			output = append(output, s)
		}
	}
	if len(songList)-len(output) > 0 {
		s.Logger.Info(fmt.Sprintf("filtered out %d songs", len(songList)-len(output)))
	}
	return output, nil
}
func (s *SongParser) DownloadSongList(songs []*Song, workers int, path string) {
	wg := sync.WaitGroup{}
	songChan := make(chan *Song, 100)
	for i := 0; i < workers; i++ {
		go func() {
			for song := range songChan {
				song.Download(s.Process, path)
				wg.Done()
			}
		}()
	}
	for _, song := range songs {
		songChan <- song
		wg.Add(1)
	}

	wg.Wait()

}
