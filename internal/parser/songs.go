package parser

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/Seann-Moser/WebParser/website"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Song struct {
	ID           string `json:"id"`
	Url          string `json:"song_url"`
	Title        string `json:"title"`
	Difficulties []string
	Thumbnail    string `json:"thumbnail"`
	DownloadLink string `json:"download_link"`
}

func (s *Song) GetID() {
	if len(s.Url) > 0 {
		sp := strings.Split(s.Url, "/")
		s.ID = sp[len(sp)-2]
	}
}
func (s *Song) GetDownloadLink() {
	if s.ID == "" {
		s.GetID()
	}
	s.DownloadLink = fmt.Sprintf("https://api.beatsaver.com/download/key/%s", s.ID)
}
func (s *Song) Download(p *ParserProcessor, path string) {
	if s.DownloadLink == "" {
		s.GetDownloadLink()
	}
	var err error
	var filename string
	files, err := ioutil.ReadDir(path)
	if err != nil && !strings.Contains(err.Error(), "system cannot find") {
		p.Logger.Error("failed reading files in dir", zap.Error(err))
		return
	}

	for _, file := range files {
		if strings.Contains(file.Name(), s.ID) {
			p.Logger.Info("file already exists skipping")
			return
		}
	}
	for i := 0; i < 3; i++ {
		filename, err = p.Req.Download(s.DownloadLink, path)
		if err == nil {
			p.Logger.Info(fmt.Sprintf("successfully downloaded %s to %s", s.Title, path))
			_, err = Unzip(filename, strings.TrimSuffix(filename, ".zip"))
			if err != nil {
				p.Logger.Error("failed unziping file "+filename, zap.Error(err))
			} else {
				os.Remove(filename)
			}

			return
		}
		time.Sleep(time.Duration(math.Pow(3, float64(i+1))) * time.Second)
	}
	if err != nil {
		p.Logger.Error(fmt.Sprintf("failed downloading %s to %s", s.Title, path), zap.Error(err))
	}

}
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
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
				Tag:   "figure",
				Order: 0,
			},
			{
				Type:     website.TypeAttribute,
				Tag:      "class",
				TagValue: "post-gallery",
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
				Tag:             "href",
				InternalTagName: "song_url",
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

func (s *SongParser) GetTopSongs(amount int) []*Song {

	return nil
}
func (s *SongParser) GetSongsWithPage(u string, amount int) ([]*Song, error) {
	var songList []*Song
	currentURL := u
	visitedMap := map[string]bool{}
	for {
		visitedMap[currentURL] = true
		sl, err := s.GetSongs(currentURL)
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
			for _, i := range p {
				if _, found := visitedMap[i.Link]; !found && len(i.Text) > 0 && i.Text != "1" {
					currentURL = i.Link
					break
				}
			}

		}
	}

}
func (s *SongParser) GetSongs(u string) ([]*Song, error) {
	b, err := s.Process.GetData(u, s.Data)
	if err != nil {
		return nil, err
	}
	var songList []*Song
	err = json.Unmarshal(b, &songList)
	if err != nil {
		return nil, err
	}
	for _, s := range songList {
		s.GetDownloadLink()
	}
	return songList, nil
}
