package utils

import (
	"reflect"
	"testing"
)

func TestUtils_SplitLines(t *testing.T) {
	type args struct {
		filepath string
	}
	tests := []struct {
		name    string
		u       *Utils
		args    args
		want    []string
		wantErr bool
	}{
		{"no file", &Utils{}, args{"fakefile"}, nil, true},
		{"empty file", &Utils{}, args{"testdata/empty"}, []string{}, false},
		{"text in a single line", &Utils{}, args{"testdata/nolines"}, []string{"everything is in one line but watch for the xtra space in the end"},
			false},
		{"some files", &Utils{}, args{"testdata/a.a"}, []string{"Some","text","on","different","lines"}, false},
		{"some files", &Utils{}, args{"testdata/mixed.txt"}, []string{"Here is an example of \\n",
			"Here is a carriage return \\n","","Some extra spaces"}, false},
		{"longline - 582KB", &Utils{}, args{"testdata/longline"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Utils{}
			got, err := u.SplitLines(tt.args.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Utils.SplitLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Utils.SplitLines() = %v, want %v", got, tt.want)
			}
		})
	}
}

