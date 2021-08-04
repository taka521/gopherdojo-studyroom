package game

import (
	"bufio"
	"context"
	"fmt"
	"io"
)

func GetTypingGame(ir io.Reader, iw io.Writer, wb WordBook) Game {
	return &typingGame{
		ir:    ir,
		iw:    iw,
		wb:    wb,
		word:  wb.Get(),
		score: &Score{},
	}
}

type typingGame struct {
	ir    io.Reader
	iw    io.Writer
	wb    WordBook
	word  Word
	score *Score
}

func (t *typingGame) Run(ctx context.Context) <-chan *Score {
	ch := make(chan *Score, 0)

	s := bufio.NewScanner(t.ir)
	go func() {
		defer close(ch)

		t.PrintWord()
		for s.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				ch <- t.Judge(s.Text()).Next().PrintWord().score
			}
		}
	}()

	return ch
}

// Next は次の単語を設定します。
func (t *typingGame) Next() *typingGame {
	t.word = t.wb.Get()
	return t
}

// PrintWord は出力先へ単語を出力します。
func (t *typingGame) PrintWord() *typingGame {
	_, _ = fmt.Fprintf(t.iw, "%s: ", t.word)
	return t
}

// Judge は指定された文字列が単語を一致するかを判定し、スコアへ反映します。
func (t *typingGame) Judge(text string) *typingGame {
	if t.word.Equals(text) {
		t.score.OK++
	} else {
		t.score.NG++
	}
	return t
}
