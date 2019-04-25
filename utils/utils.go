package utils

import (
	"bufio"
	"os"
)

type Utilizer interface {
	Filer
}

type Filer interface {
	SplitLines(filepath string) ([]string, error)
}

type Utils struct {}

// SplitLines takes a filepath, reads the file and returns all the lines as a string array.
// An existent empty file will return nil
// Error can only be from not finding the file, or some issue during the read.
func (u *Utils) SplitLines(filepath string) ([]string, error) {
	lines := make([]string, 0)
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}