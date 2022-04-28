package parser

import (
	"archive/zip"
	"fmt"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Song struct {
	ID              string `json:"id"`
	Url             string `json:"song_url"`
	Title           string `json:"title"`
	Difficulties    []string
	Thumbnail       string `json:"thumbnail"`
	DownloadLink    string `json:"download_link"`
	RawDifficulties string `json:"raw_difficulties"`
	RatingPercent   float32
	RawText         string `json:"raw_text"`
}

func (s *Song) Process() {
	s.GetID()
	s.GetDownloadLink()
	s.GetTitle()
	s.ProcessRawText()
}

func (s *Song) ProcessRawText() {
	if strings.Contains(s.RawText, "---") {
		sp := strings.Split(s.RawText, "---")
		titleIndex := -1
		for i, v := range sp {
			if strings.EqualFold(v, s.Title) {
				titleIndex = i
				break
			}
		}
		if titleIndex == -1 {
			return
		}
		difficultyIndex := titleIndex + 1
		difficulties := sp[difficultyIndex:]
		rating := []string{}
		for i, d := range difficulties {
			if _, err := strconv.Atoi(d); err == nil {
				rating = sp[difficultyIndex+i : difficultyIndex+i+2]
				break
			} else if strings.EqualFold(d, "difficulties") {
				continue
			} else {
				s.Difficulties = append(s.Difficulties, d)
			}
		}
		thumbsUp, err := strconv.Atoi(rating[0])
		if err != nil {

		}
		thumbsDown, err := strconv.Atoi(rating[1])
		if err != nil {

		}
		if thumbsDown+thumbsUp > 0 {
			p := float32(thumbsUp) / float32(thumbsUp+thumbsDown)
			s.RatingPercent = p
		}

	}
}

func (s *Song) GetTitle() {
	if strings.Contains(s.Title, ",") {
		sp := strings.Split(s.Title, ",")
		s.Title = sp[0]
	}
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
func (s *Song) AlreadyDownloaded(path string) bool {
	files, err := ioutil.ReadDir(path)
	if err != nil && !strings.Contains(err.Error(), "system cannot find") {
		return false
	}

	for _, file := range files {
		if strings.Contains(file.Name(), s.ID) {
			return true
		}
	}
	return false
}
func (s *Song) Download(p *ParserProcessor, path string) {
	if s.DownloadLink == "" {
		s.GetDownloadLink()
	}
	var err error
	var filename string
	if s.AlreadyDownloaded(path) {
		p.Logger.Debug("file already exists skipping")
		return
	}
	for i := 0; i < 3; i++ {
		filename, err = p.Req.Download(s.DownloadLink, path)
		if err == nil {
			p.Logger.Info(fmt.Sprintf("successfully downloaded %s to %s", s.Title, path))
			_, err = Unzip(filename, strings.TrimSuffix(filename, ".zip"))
			if err != nil {
				p.Logger.Error("failed unzipping file "+filename, zap.Error(err))
			}
			_ = os.Remove(filename)
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
