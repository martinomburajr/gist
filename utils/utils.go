package utils

import (
	"fmt"
	"github.com/martinomburajr/gist/gists"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//SendAllGistFiles sends a set of gist files to Github simultaneously
func SendAllGistFiles(gistss []*gists.GistFile) (chan *gists.GistFile, chan *http.Response) {
	failedChan := make(chan *gists.GistFile, 0)
	responses := make(chan *http.Response, 0)
	for _, v := range gistss {
		go func(gist *gists.GistFile) {
			response, err := gist.Create()
			if err != nil {
				_, err = io.Copy(os.Stderr, strings.NewReader(err.Error()))
				failedChan <- gist
			} else {
				responses <- response
			}
		}(v)
	}
	return failedChan, responses
}

//ScanAllFilesInDir checks to see if files in a given directory are gistable and resturns only the gistable ones
func ScanAllFilesInDir(dir string) []*gists.GistFile {
	filepaths := GetAllFilesInDir(dir)
	gistfiles := make([]*gists.GistFile, 0)
	for _, v := range filepaths {
		gist := gists.GistParser{
			Filepath: v,
		}
		gistFile, err := gist.ToGist()
		if err != nil {

		} else {
			gistfiles = append(gistfiles, gistFile)
		}
	}
	return gistfiles
}

//GetAllFilesInDir returns a list of files in a given directory
func GetAllFilesInDir(dir string) []string {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
	return files
}