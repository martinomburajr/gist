package gists

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var DummyGistFile1 =  GistFile{
	Description: "",
	Files: []GistFileBody{ {Content: "test-a.a"}},
	Public: false,
}


func TestGistFile_Delete(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		g       *GistFile
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.Delete(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GistFile.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GistFile.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGistFile_Update(t *testing.T) {
	type args struct {
		in0 interface{}
	}
	tests := []struct {
		name    string
		g       *GistFile
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.Update(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("GistFile.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GistFile.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGistFile_Create(t *testing.T) {
	data, err := json.Marshal(DummyGistFile1)
	if err != nil {
		t.Error(err)
	}

	mux := http.DefaultServeMux
	req, err := http.NewRequest(http.MethodPost, "/gists",  bytes.NewReader(data))
	if err != nil {
		log.Fatal("")
	}
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, req)

	if status := recorder.Code; recorder.Code > 399 {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	//tests := []struct {
	//	name    string
	//	g       *GistFile
	//	want    *http.Response
	//	wantErr bool
	//}{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		got, err := tt.g.Create()
	//		if (err != nil) != tt.wantErr {
	//			t.Errorf("GistFile.Create() error = %v, wantErr %v", err, tt.wantErr)
	//			return
	//		}
	//		if !reflect.DeepEqual(got, tt.want) {
	//			t.Errorf("GistFile.Create() = %v, want %v", got, tt.want)
	//		}
	//	})
	//}
}

func TestGistFile_Retrieve(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		g       *GistFile
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.g.Retrieve(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GistFile.Retrieve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GistFile.Retrieve() = %v, want %v", got, tt.want)
			}
		})
	}
}
