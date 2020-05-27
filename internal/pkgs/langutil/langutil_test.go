package langutil

import (
	"reflect"
	"testing"

	"golang.org/x/text/language"
)

func Test_lang(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want language.Tag
	}{
		{
			name: "browser-firefox",
			args: args{
				s: "en-US,en;q=0.9,tr-TR;q=0.8,tr;q=0.7",
			},
			want: language.English,
		},
		{
			name: "empty",
			args: args{
				s: "",
			},
			want: language.English,
		},
		{
			name: "malformed",
			args: args{
				s: "afsjg;lkafgslasnl",
			},
			want: language.English,
		},
		{
			name: "tr",
			args: args{
				s: "tr-TR;q=0.2",
			},
			want: language.Turkish,
		},
		{
			name: "en",
			args: args{
				s: "en-US,tr-TR;q=0.9",
			},
			want: language.English,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LangFromAcceptLanguage(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lang() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchable(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "trlower",
			text: "qwertyuıopğü,asdfghjklşizxcvbnmöç.",
			want: "QWERTYUIOPGUASDFGHJKLSIZXCVBNMOC",
		},
		{
			name: "trUPPER",
			text: "qwertyuıopğü,asdfghjklşizxcvbnmöç.",
			want: "QWERTYUIOPGUASDFGHJKLSIZXCVBNMOC",
		},
		{
			name: "punct",
			text: "b!@#$%^&*()_+~ !~!@#$%^&*()_+e",
			want: "BE",
		},
		{
			name: "malformed",
			// malformed input should ignore that field incase there is a malformed data sent from client
			text: string([]byte{'a', 0xFF, 'a'}),
			want: "AA",
		},
		{
			name: "bei jing",
			text: "北京kožušček",
			want: "BEIJINGKOZUSCEK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Searchable(tt.text); got != tt.want {
				t.Errorf("Searchable() = %v, want %v", got, tt.want)
			}
		})
	}
}
