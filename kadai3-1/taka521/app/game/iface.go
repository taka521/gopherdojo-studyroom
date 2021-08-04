//go:generate gomockhandler -config=../../gomockhandler.json -source=$GOFILE -destination=./mock/mock_$GOFILE -package=mock

package game

import "context"

// Score はゲームのスコアを表します
type Score struct {
	OK uint // 正解数
	NG uint // 不正解数
}

type Game interface {
	// Run はゲームを開始します
	Run(ctx context.Context) <-chan *Score
}
