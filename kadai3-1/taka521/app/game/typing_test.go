package game

import (
	"bytes"
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/taka521/gopherdojo-studyroom/kadai3-1/taka521/app/game/mock"
)

func TestGetTypingGame(t *testing.T) {
	ctrl := gomock.NewController(t)

	ir := mock.NewMockReader(ctrl)
	iw := mock.NewMockWriter(ctrl)
	wb := NewMockWordBook(ctrl)
	wb.EXPECT().Get().Return(Word("test"))

	type args struct {
		ir io.Reader
		iw io.Writer
		wb WordBook
	}
	tests := []struct {
		name  string
		args  args
		want  Game
		setup func()
	}{
		{
			name: "typingGame のポインタが返ること",
			args: args{ir: ir, iw: iw, wb: wb},
			want: &typingGame{ir: ir, iw: iw, wb: wb, word: "test", score: &Score{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTypingGame(tt.args.ir, iw, tt.args.wb)
			opt := cmpopts.IgnoreFields(typingGame{}, "ir", "iw", "wb", "word", "score")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("GetTypingGame() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_typingGame_Judge(t *testing.T) {
	type fields struct {
		ir    io.Reader
		iw    io.Writer
		wb    WordBook
		word  Word
		score *Score
	}
	type args struct {
		text string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Score
	}{
		{
			name:   "指定された text が word と一致する場合、Score.OK に加算されること",
			fields: fields{word: "test", score: &Score{}},
			args:   args{text: "test"},
			want:   &Score{OK: 1, NG: 0},
		},
		{
			name:   "指定された text が word と一致しない場合、Score.NG に加算されること",
			fields: fields{word: "test", score: &Score{}},
			args:   args{text: "hoge"},
			want:   &Score{OK: 0, NG: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			g := &typingGame{
				ir:    tt.fields.ir,
				iw:    tt.fields.iw,
				wb:    tt.fields.wb,
				word:  tt.fields.word,
				score: tt.fields.score,
			}
			got := g.Judge(tt.args.text).score
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetTypingGame() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_typingGame_Next(t1 *testing.T) {
	ctrl := gomock.NewController(t1)

	wb := NewMockWordBook(ctrl)
	wb.EXPECT().Get().Return(Word("test"))

	type fields struct {
		ir    io.Reader
		iw    io.Writer
		wb    WordBook
		word  Word
		score *Score
	}
	tests := []struct {
		name   string
		fields fields
		want   Word
	}{
		{
			name:   "WorkBook.Get() の値が word に設定されること",
			fields: fields{wb: wb, word: "golang"},
			want:   "test",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &typingGame{
				ir:    tt.fields.ir,
				iw:    tt.fields.iw,
				wb:    tt.fields.wb,
				word:  tt.fields.word,
				score: tt.fields.score,
			}
			if got := t.Next().word; !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_typingGame_PrintWord(t1 *testing.T) {
	ctrl := gomock.NewController(t1)

	iw := mock.NewMockWriter(ctrl)

	type fields struct {
		ir    io.Reader
		iw    io.Writer
		wb    WordBook
		word  Word
		score *Score
	}
	tests := []struct {
		name   string
		fields fields
		want   func()
	}{
		{
			name:   "io.Writer に期待する文字列が渡されること",
			fields: fields{iw: iw, word: "test"},
			want: func() {
				iw.EXPECT().Write([]byte("test: ")).Return(len("test: "), nil)
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			tt.want()
			t := &typingGame{
				ir:    tt.fields.ir,
				iw:    tt.fields.iw,
				wb:    tt.fields.wb,
				word:  tt.fields.word,
				score: tt.fields.score,
			}
			t.PrintWord()
		})
	}
}

func Test_typingGame_Run(t1 *testing.T) {
	ctrl := gomock.NewController(t1)

	ir := bytes.NewBufferString("one\nthree\ntwo")

	iw := mock.NewMockWriter(ctrl)
	iw.EXPECT().Write(gomock.Any()).Return(0, nil).Times(3)

	wb := NewMockWordBook(ctrl)
	wb.EXPECT().Get().Return(Word("two"))
	wb.EXPECT().Get().Return(Word("three"))

	type fields struct {
		ir    io.Reader
		iw    io.Writer
		wb    WordBook
		word  Word
		score *Score
	}
	tests := []struct {
		name   string
		fields fields
		want   *Score
	}{
		{
			name:   "context がキャンセルされるまでは、処理が実行されること",
			fields: fields{ir: ir, iw: iw, wb: wb, word: "one", score: &Score{}},
			want:   &Score{OK: 1, NG: 1},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &typingGame{
				ir:    tt.fields.ir,
				iw:    tt.fields.iw,
				wb:    tt.fields.wb,
				word:  tt.fields.word,
				score: tt.fields.score,
			}

			ctx, cancel := context.WithCancel(context.Background())
			var got *Score
			ch := t.Run(ctx)

			got = <-ch // 1回目の受信
			cancel()
			got = <-ch // 2回目の受信 -> 3回目のループで ctx.Done のケースに入る

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t1.Errorf("Run() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
