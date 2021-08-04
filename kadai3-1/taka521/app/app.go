package app

import (
	"context"
	"fmt"
	"time"

	"github.com/taka521/gopherdojo-studyroom/kadai3-1/taka521/app/game"
)

func App(ctx context.Context, limit uint, g game.Game) string {
	cwt, cancel := context.WithTimeout(ctx, time.Duration(limit)*time.Second)
	defer cancel()

	gc := g.Run(cwt)
	var s *game.Score
	for {
		select {
		case <-cwt.Done():
			return fmt.Sprintf("\n==== Score: %v/%v ====", s.OK, s.OK+s.NG)
		case s = <-gc:
		}
	}
}
