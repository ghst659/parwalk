package parwalk

import (
	"context"
	"fmt"
	"os"
	"testing"
)

type datum struct {
	path string
	size int64
}

func TestParWalk(t *testing.T) {
	const root = "/opt"
	data := make(chan datum)
	go func(ctx context.Context, top string) {
		defer close(data)
		act := func(ctx context.Context, pth string) bool {
			st, err := os.Stat(pth)
			if err != nil {
				return false
			}
			data <- datum{
				path: pth,
				size: st.Size(),
			}
			return true
		}
		Traverse(ctx, top, act)
	}(context.Background(), root)
	var total int64
	for d := range data {
		fmt.Printf("  %s: %d\n", d.path, d.size)
		total += d.size
	}
	fmt.Printf("Total: %d\n", total)
}
