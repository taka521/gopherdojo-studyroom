package main

import (
	"context"

	"github.com/taka521/gopherdojo-studyroom/kadai3-1/taka521/app"
)

func main() {
	ctx := context.Background()
	app.App(ctx, 30)
}
