package app

import (
	"bufio"
	"context"
	"fmt"
	"os"
)

type result struct {
	ok int // 正解数
	ng int // 不正解数
}

func typing(ctx context.Context) <-chan *result {
	ch := make(chan *result, 0)
	word := GetWord()
	r := &result{}

	s := bufio.NewScanner(os.Stdin)
	go func() {
		defer close(ch)

		fmt.Printf("%s: ", word)
		for s.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				if word.Equals(s.Text()) {
					r.ok++
				} else {
					r.ng++
				}
				word = GetWord()
				fmt.Printf("%s: ", word)
				ch <- r
			}
		}
	}()

	return ch
}
