package main

import (
	"context"
	"fmt"
	"os"

	"github.com/taka521/gopherdojo-studyroom/kadai3-1/taka521/app"
	"github.com/taka521/gopherdojo-studyroom/kadai3-1/taka521/app/game"
)

func main() {
	ctx := context.Background()
	typing := game.GetTypingGame(os.Stdin, os.Stdout, game.GetWordBook())
	fmt.Println(app.App(ctx, 30, typing))
}
