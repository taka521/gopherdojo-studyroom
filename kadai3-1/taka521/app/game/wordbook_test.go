package game

import (
	"reflect"
	"testing"
)

func TestGetWordBook(t *testing.T) {
	tests := []struct {
		name string
		want WordBook
	}{
		{
			name: "wordbook のポインタが返ること",
			want: &wordbook{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetWordBook(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWordBook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWord_Equals(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		w    Word
		args args
		want bool
	}{
		{
			name: "文字列が全て一致する場合 true が返ること",
			w:    "test",
			args: args{v: "test"},
			want: true,
		},
		{
			name: "大文字・小文字が異なる場合でも true が返ること",
			w:    "test",
			args: args{v: "TEST"},
			want: true,
		},
		{
			name: "文字列が一致しない場合は false が返ること",
			w:    "test",
			args: args{v: "tset"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.Equals(tt.args.v); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_wordbook_Get(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "words に存在する単語が返ること　",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &wordbook{}
			got := w.Get()

			for _, word := range words {
				if word == got {
					return
				}
			}
			t.Errorf("words に存在しない値が")
		})
	}
}
