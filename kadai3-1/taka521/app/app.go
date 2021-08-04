package app

import (
	"context"
	"fmt"
	"time"
)

func App(ctx context.Context, limit uint) {
	cwt, cancel := context.WithTimeout(ctx, time.Duration(limit)*time.Second)
	defer cancel()

	gc := typing(cwt)
	var r *result = nil
	for {
		select {
		case <-cwt.Done():
			fmt.Printf("\n==== Score: %v/%v ====", r.ok, r.ok+r.ng)
			return
		case r = <-gc:
		}
	}
}
