//go:generate gomockhandler -config=../../gomockhandler.json -source=$GOFILE -destination=iface_mock.go -package=$GOPACKAGE

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
