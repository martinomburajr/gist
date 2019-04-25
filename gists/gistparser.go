package gists

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)


type GistParser struct {
	//File path refers to the file path. It is a simple string e.g. "/home/me/files"
	Filepath string `json:"badfilepath"`

	//fileContents are all the elements within the file regardless of whether they are gistable or not
	fileContents []byte
}

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

func (g *GistParser) GetFileBody() (*GistFileBody, error) {
	if len(g.fileContents) == 0 {
		err := g.Reader()
		return nil, err
	}

	return &GistFileBody{
		Content: string(g.fileContents),
	},nil
}

// IsGistable checks to a see a certain file conforms to the gist standard.
// For a file to be considered "gistable". It must contain the text "start gist" <some code> "end gist" where
// "start gist" and "end gist"  (without quotes) legitimize the file as being gistable.
// Also note the "start gist" and "end gist" are case insensitive.
// These can be placed within comments, or in plain text so long as these two elements are there one after the other.
// the <some-code> refers to code you would like to send to gist.github.
// com gist attributes such as name can be placed within the file between the start  gist and end gist sections. e.g.
// gist name: "some gist". See the respective methods for more.
func (g *GistParser) IsGistable() error {
	//perform readfile
	err := g.Reader()
	if err != nil {
		return err
	}

	lines, err := g.getgistLines()
	containsStart := -2
	containsEnd := -1
	for i, v := range lines {
		if strings.Contains(strings.ToLower(v), "start gist") {
			containsStart = i
		}
		if strings.Contains(strings.ToLower(v), "end gist") {
			containsEnd = i
		}
	}

	if containsStart < 0 {
		return fmt.Errorf("file is not gistable. no contains start found")
	}

	if containsEnd < 0 {
		return fmt.Errorf("file is not gistable. no contains end found")
	}

	if  containsStart >= containsEnd{
		return fmt.Errorf("start gist and end gist cannot be on the same line. " +
			"end gist must come at least 1 line after start gist. There must be some code content in between (" +
			"excluding metadata e.g. gist name: or gist description etc)")
	}

	return nil
}

// getAuthor returns the Author information. The author must precede the gist label, and must contain all runes within the word
// "AUTHOR". This is CASE sensitive
//
// This is the format shown below. Email is optional. Anything after newline carriage return is considered not part of the author label
//
//	/** Start gist
//	Author: I am some author <hereismy@email.com>
//  Description: Some awesome gist
//  Public: true
//  end gist
//	*/
//  returns I am some author <hereismy@email.com>
//
func (g *GistParser) GetAuthor() (string, error) {
	lines, err := g.getgistLines()
	if err != nil {
		return "", err
	}
	return g.getContent(lines, "author")
}

// getDescription returns the Author information. The Description must precede the gist label, and must contain all runes within the word
// "Description". This is CASE sensitive
//
// This is the format shown below. Email is optional. Anything after newline carriage return is considered not part of the Description label so ensure description is all in a single line.
//
//	/** Start gist
//	Author: I am some author <hereismy@email.com>
//  Description: Some awesome gist
//  Public: true
//  end gist
//	*/
//  returns Some awesome gist
//
func (g *GistParser) GetDescription() (string, error) {
	lines, err := g.getgistLines()
	if err != nil {
		return "", err
	}
	return g.getContent(lines, "description")
}

// getPublic returns the isPublic information. The Public must precede the gist label, and must contain all runes within the word
// "AUTHOR". This is CASE insensitive
//
// This is the format shown below. Email is optional. Public is a boolean variable that can either be true or false. Its default value is true
//
//	/** Start gist
//	Author: I am some author <hereismy@email.com>
//  Description: Some awesome gist
//  Public: true
//  end gist
//	*/
// returns true
//
func (g *GistParser) GetPublic() (bool, error) {
	lines, err := g.getgistLines()
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

//getgistLines returns the lines encapsulated by the 'start gist' and the'end gist' labels. This is where all the important gist metadata is found.
func (g *GistParser) getgistLines() ([]string, error) {
	err := g.IsGistable()
	if err != nil {
		return nil, fmt.Errorf("could not determine if file has 'start gist' and the'end gist' labels -> %s", err.Error())
	}

	buffer := bytes.NewBuffer(g.fileContents)
	documentContents := strings.SplitAfter(buffer.String(), "\n")

	startIndex := -1
	endIndex := -1
	for i, _ := range documentContents {
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

//getContent takes in the gist section obtained after running getgistLines, and obtaining the exact metadata section. key represents a  key in a key-value pair. e.g. author or description are valid keys
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
	}  else {
		return "", fmt.Errorf(key + " does not exist")
	}
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

//Reads the entire file contents using ioutil. It then appends the data to the p variable for later use.
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
