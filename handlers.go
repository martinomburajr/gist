package main

import (
	"fmt"
	"github.com/martinomburajr/gist/gists"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

/**
	start GIST
	Author: "Martin Ombura Jr." <info@martinomburajr.com>"
	Description: ""
 */

func CreateGistHandler(w http.ResponseWriter, r *http.Request) {

}

func SendAllGistFiles(gistss []*gists.GistFile) (chan *gists.GistFile, chan *http.Response) {
	failedChan := make(chan *gists.GistFile, 0)
	responses := make(chan *http.Response, 0)
	for _, v := range gistss {
		go func(gist *gists.GistFile) {
			response, err := gist.Create()
			if err != nil {
				_, err = io.Copy(os.Stderr, strings.NewReader(err.Error()))
				failedChan <- gist
			}else {
				responses <- response
			}
		}(v)
	}
	return failedChan, responses
}

func ScanAllFilesInDir(dir string) []*gists.GistFile {
	filepaths := getAllFilesInDir(dir)
	gistfiles := make([]*gists.GistFile,0)
	for _, v := range filepaths {
		gist := gists.GistParser{
			Filepath: v,
		}
		gistFile, err := gist.ToGist()
		if err != nil {

		}else {
			gistfiles = append(gistfiles, gistFile)
		}
	}
	return gistfiles
}

func getAllFilesInDir(dir string) []string {
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