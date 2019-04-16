package gists

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

//GistParser is an object that represents the items required to parse a gist. Typically the filepath and its contents
type GistParser struct {
	Filepath string `json:"filepath"`
	fileContents []byte
}

//ToGist is an accumulator method that performs multiple sub functions.
// It is responsible for the creation of a gist file, this involves, parsing the  file to determine if it s gistable,
// if it is gistable, relavant metadata is extracted from the file, the file is then converted to the GistFile type.
// This is the final form before other operations can be performed on it once it is a GistFile.
func (g *GistParser) ToGist() (*GistFile, error) {
	err := g.IsGistable()
	if err != nil {
		return nil, fmt.Errorf("could not get body -> %s", err.Error())
	}

	description, err := g.GetDescription()
	if err != nil {
		return nil, err
	}

	if len(g.fileContents) == 0 {
		err := g.Reader()
		return nil, err
	}

	b, err := g.GetPublic()
	if err != nil {
		return nil, err
	}

	gistFileBody, err := g.GetFileBody()
	if err != nil {
		return nil, err
	}

	files := []GistFileBody{*gistFileBody}
	return &GistFile{
		Description: description,
		Files: files,
		Public: b,
	}, nil
}

//GetFileBody extracts the body of a gist from the file.
func (g *GistParser) GetFileBody() (*GistFileBody, error) {
	if len(g.fileContents) == 0 {
		err := g.Reader()
		return nil, err
	}

	return &GistFileBody{
		Content: string(g.fileContents),
	},nil
}

// IsGistable checks to a see a certain file conforms to the GOGIST standard.
// If the file does not contain the "GOGIST" label
// in a comments section in the file. It is deemed ungistable meaning, gist will not create a gist for the user in that regard.
// If no error is presented, one can assume that it is a gistable file.
func (g *GistParser) IsGistable() error {
	err := g.Reader()
	if err != nil {
		return err
	}

	containsStart := bytes.Contains(bytes.ToLower(g.fileContents), bytes.ToLower([]byte("start GOGIST")))
	if !containsStart {
		return fmt.Errorf("is not a suitable GOGIST file. Add the following string 'start GOGIST' inside a comment section at the top of the file to mark it as a file gist can gist ;-)")
	}

	containsEnd := bytes.Contains(bytes.ToLower(g.fileContents), bytes.ToLower([]byte("end GOGIST")))
	if !containsEnd {
		return fmt.Errorf("is not a suitable GOGIST file. Add the following string 'end GOGIST' inside a comment section at the top of the file to mark it as a file gist can gist ;-). This should be after the start Gogist section")
	}

	return nil
}

// GetAuthor returns the Author information. The author must precede the GOGIST label, and must contain all runes within the word
// "AUTHOR". This is CASE sensitive
//
// This is the format shown below. Email is optional. Anything after newline carriage return is considered not part of the author label
//
//	/** Start GOGIST
//	Author: I am some author <hereismy@email.com>
//  Description: Some awesome gist
//  Public: true
//  end gist
//	*/
//  returns I am some author <hereismy@email.com>
//
func (g *GistParser) GetAuthor() (string, error) {
	lines, err := g.getGogistLines()
	if err != nil {
		return "", err
	}
	return g.getContent(lines, "author")
}

// GetDescription returns the Author information. The Description must precede the GOGIST label,
// and must contain all runes within the word
// "Description". This is CASE sensitive
//
// This is the format shown below. Email is optional. Anything after newline carriage return is considered not part of the Description label so ensure description is all in a single line.
//
//	/** Start GOGIST
//	Author: I am some author <hereismy@email.com>
//  Description: Some awesome gist
//  Public: true
//  end gist
//	*/
//  returns Some awesome gist
//
func (g *GistParser) GetDescription() (string, error) {
	lines, err := g.getGogistLines()
	if err != nil {
		return "", err
	}
	return g.getContent(lines, "description")
}

// GetPublic returns the isPublic information. The Public must precede the GOGIST label,
// and must contain all runes within the word
// "AUTHOR". This is CASE insensitive
//
// This is the format shown below. Email is optional. Public is a boolean variable that can either be true or false. Its default value is true
//
//	/** Start GOGIST
//	Author: I am some author <hereismy@email.com>
//  Description: Some awesome gist
//  Public: true
//  end gist
//	*/
// returns true
//
func (g *GistParser) GetPublic() (bool, error) {
	lines, err := g.getGogistLines()
	if err != nil {
		return false, err
	}
	content, err := g.getContent(lines, "public")
	if err != nil {
		return true, nil
	}
	b, err := strconv.ParseBool(content)
	if err != nil {
		return false, fmt.Errorf("couldnt parse string to bool -> %v", err)
	}
	return b, nil
}

//getGogistLines returns the lines encapsulated by the 'start gist' and the'end gist' labels. This is where all the important gist metadata is found.
func (g *GistParser) getGogistLines() ([]string, error) {
	err := g.IsGistable()
	if err != nil {
		return nil, fmt.Errorf("could not determine if file has 'start gist' and the'end gist' labels -> %s", err.Error())
	}

	buffer := bytes.NewBuffer(g.fileContents)
	documentContents := strings.SplitAfter(buffer.String(), "\n")

	startIndex := -1
	endIndex := -1
	for i := range documentContents {
		documentContents[i] = strings.Trim(documentContents[i], " \r\n")
		if strings.Contains(  strings.ToLower(documentContents[i]), strings.ToLower("start gist"))  {
			startIndex = i
		}
		if strings.Contains(strings.ToLower(documentContents[i]), strings.ToLower("end gist"))  {
			endIndex = i
			break
		}
	}

	return documentContents[startIndex:endIndex+1], nil
}

//getContent takes in the gist section obtained after running getGogistLines, and obtaining the exact metadata section. key represents a  key in a key-value pair. e.g. author or description are valid keys
func (g *GistParser) getContent(s []string, key string) (string, error) {
	var location = -1
	for i, v := range s {
		if strings.Contains(strings.ToLower(v), strings.ToLower(key)) {
			location = i
			break
		}
	}
	if location > 0 {
		totalString := strings.TrimSpace(s[location])
		replaced := strings.Replace(strings.ToLower(totalString), strings.ToLower(key)+":", "", -1)
		lenreplaced := len(replaced)
		finstring := totalString[(len(totalString) - lenreplaced):]
		trim := strings.Trim(finstring, " ")
		return trim, nil
	}
	return "", fmt.Errorf(key + " does not exist")
}


//Reads the entire file contents using ioutil. It then appends the data to the p variable for later use.
//func (g *GistParser) Read(p []byte) (n int, err error) {
//	data, err := ioutil.ReadFile(g.Filepath)
//	if err != nil {
//		return 0, fmt.Errorf("file may not exist -> %s", err)
//	}
//	for _, v := range data {
//		p = append(p, v)
//	}
//	return len(data), nil
//}

//Reader Reads the entire file contents using ioutil. It then appends the data to the p variable for later use.
func (g *GistParser) Reader()  (err error) {
	data, err := ioutil.ReadFile(g.Filepath)
	if err != nil {
		return  fmt.Errorf("file may not exist -> %s", err)
	}
	g.fileContents = data
	return nil
}

type gistparsereader struct {
	done bool
	read []byte
}
