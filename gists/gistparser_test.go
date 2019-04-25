package gists

import (
	"io/ioutil"
	"reflect"
	"testing"
)

//.c file extensions are typically problematic when kept in the testdata folder.
// todo place gistparser_test initialization variables in a different folder same package
var badfilepath = "samplefile.randextensions"
var filepatha = "testdata/test-a.a"
var filepathb = "testdata/test-b.b"
var filepathc = "testdata/test-c.d"
var filepathe = "testdata/test-e.e"
var filepathf = "testdata/test-f.f"
var fielpathg = "testdata/test-g.g"
var filepathgo = "testdata/test-go.go"
var filepathrandom = "testdata/test-random.rand"

var gistsectiona, _ = (&GistParser{filepatha, nil}).getgistLines()
var gistsectionb, _ = (&GistParser{filepathb, nil}).getgistLines()
var gistsectionc, _ = (&GistParser{filepathc, nil}).getgistLines()
var gistsectiongo, _ = (&GistParser{filepathgo, nil}).getgistLines()
var gistsectionrand, _ = (&GistParser{filepathrandom, nil}).getgistLines()

func readFile(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func TestGistParser_ToGist(t *testing.T) {
	type fields struct {
		Filepath     string
		fileContents []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    *GistFile
		wantErr bool
	}{
		{"bad-file", fields{Filepath: badfilepath, fileContents: nil}, nil, true},
		{"samplefile-a", fields{Filepath: filepatha, fileContents: nil}, nil, true},
		//This will only get the first line as the attributes should not be multilined
		{"samplefile-b", fields{Filepath: filepathb, fileContents: nil}, &GistFile{
			Description: `the following program will calculate the constant e-2 to about`,
			Public:      true,
			Files:       []GistFileBody{{readFile(filepathb)}},
		}, false},
		{"samplefile-c", fields{Filepath: filepathc, fileContents: nil}, &GistFile{
			Description: `How to create random vars in C`,
			Public:      true,
			Files:       []GistFileBody{{readFile(filepathc)}},
		}, false},
		{"samplefile-go", fields{Filepath: filepathgo, fileContents: nil}, &GistFile{
			Description: `How to create a server in Go`,
			Public:      false,
			Files:       []GistFileBody{{readFile(filepathgo)}},
		}, false},
		{"samplefile-random", fields{Filepath: filepathrandom, fileContents: nil}, &GistFile{
			Description: "_fnsofld",
			Public:      true,
			Files:       []GistFileBody{{readFile(filepathrandom)}},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GistParser{
				Filepath:     tt.fields.Filepath,
				fileContents: tt.fields.fileContents,
			}
			got, err := g.ToGist()
			if (err != nil) != tt.wantErr {
				t.Errorf("GistParser.ToGist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GistParser.ToGist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGistParser_GetFileBody(t *testing.T) {
	type fields struct {
		Filepath     string
		fileContents []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    *GistFileBody
		wantErr bool
	}{
		//{"", }
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GistParser{
				Filepath:     tt.fields.Filepath,
				fileContents: tt.fields.fileContents,
			}
			got, err := g.GetFileBody()
			if (err != nil) != tt.wantErr {
				t.Errorf("GistParser.GetFileBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GistParser.GetFileBody() = %v, want %v", got, tt.want)
			}
		})
	}
}



func TestGistParser_IsGistable(t *testing.T) {
	type fields struct {
		Filepath     string
		fileContents []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"bad-file-path", fields{Filepath: badfilepath, fileContents: nil}, true},
		{"test-a - no end to gist", fields{Filepath: filepatha, fileContents: nil}, true},
		{"test-b - gist", fields{Filepath: filepathb, fileContents: nil}, false},
		{"test-c - gist", fields{Filepath: filepathc, fileContents: nil}, false},
		{"test-e - gist in between code and file", fields{Filepath: filepathe,
			fileContents: nil}, false},
		{"test-f - gist missing start but contains end", fields{Filepath: filepathf,
			fileContents: nil}, true},
		{"test-g - end gist swapped with start gist", fields{Filepath: fielpathg,
			fileContents: nil}, true},
		{"test-go - gist", fields{Filepath: filepathgo, fileContents: nil}, false},
		{"test-random - gist", fields{Filepath: filepathrandom, fileContents: nil},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GistParser{
				Filepath:     tt.fields.Filepath,
				fileContents: tt.fields.fileContents,
			}
			if err := g.IsGistable(); (err != nil) != tt.wantErr {
				t.Errorf("GistParser.IsGistable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGistParser_GetAuthor(t *testing.T) {
	type fields struct {
		Filepath     string
		fileContents []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GistParser{
				Filepath:     tt.fields.Filepath,
				fileContents: tt.fields.fileContents,
			}
			got, err := g.GetAuthor()
			if (err != nil) != tt.wantErr {
				t.Errorf("GistParser.GetAuthor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GistParser.GetAuthor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGistParser_GetDescription(t *testing.T) {
	type fields struct {
		Filepath     string
		fileContents []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GistParser{
				Filepath:     tt.fields.Filepath,
				fileContents: tt.fields.fileContents,
			}
			got, err := g.GetDescription()
			if (err != nil) != tt.wantErr {
				t.Errorf("GistParser.GetDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GistParser.GetDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGistParser_GetPublic(t *testing.T) {
	type fields struct {
		Filepath     string
		fileContents []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GistParser{
				Filepath:     tt.fields.Filepath,
				fileContents: tt.fields.fileContents,
			}
			got, err := g.GetPublic()
			if (err != nil) != tt.wantErr {
				t.Errorf("GistParser.GetPublic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GistParser.GetPublic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGistParser_getgistLines(t *testing.T) {
	type fields struct {
		Filepath     string
		fileContents []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{name: "test-file-exists", fields: fields{Filepath: badfilepath}, want: nil, wantErr: true},
		{name: "test-gister", fields: fields{Filepath: filepathrandom}, want: []string{"/* START gist",
			"Author: Martin Ombura", "Description: _fnsofld", "END gist"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GistParser{
				Filepath:     tt.fields.Filepath,
				fileContents: tt.fields.fileContents,
			}
			got, err := g.getgistLines()
			if (err != nil) != tt.wantErr {
				t.Errorf("GistParser.getgistLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GistParser.getgistLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGistParser_getContent(t *testing.T) {
	type fields struct {
		Filepath     string
		fileContents []byte
	}
	type args struct {
		s   []string
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{name: "this should recognize the lack of end gist", fields: fields{filepatha, nil}, args: args{gistsectiona,
			"description"},
			want: "", wantErr: true},
		{name: "should correctly obtain definition contents", fields: fields{filepathb, nil}, args: args{gistsectionb,
			"description"},
			want: "the following program will calculate the constant e-2 to about", wantErr: false},
		//{name: "xxx-xxx", fields:fields{filepatha, nil}, args:args{gistsectionc, "definition"},
		//	want:"This gist has no end", wantErr:false},
		//{name: "xxx-xxx", fields:fields{filepatha, nil}, args:args{gistsectiongo, "definition"},
		//	want:"This gist has no end", wantErr:false},
		//{name: "xxx-xxx", fields:fields{filepatha, nil}, args:args{gistsectionrand, "definition"},
		//	want:"This gist has no end", wantErr:false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GistParser{
				Filepath:     tt.fields.Filepath,
				fileContents: tt.fields.fileContents,
			}
			got, err := g.getContent(tt.args.s, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GistParser.getContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GistParser.getContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGistParser_Reader(t *testing.T) {
	type fields struct {
		Filepath     string
		fileContents []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GistParser{
				Filepath:     tt.fields.Filepath,
				fileContents: tt.fields.fileContents,
			}
			if err := g.Reader(); (err != nil) != tt.wantErr {
				t.Errorf("GistParser.Reader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
