// Package parwalk traverses the filesystem producing per-node
// results on a channel.
package parwalk

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type NodeAction func(context.Context, string) bool

func Traverse(ctx context.Context, root string, action NodeAction) {
	var wg sync.WaitGroup
	wg.Add(1)
	go visit(ctx, &wg, root, action)
	wg.Wait()
}

func visit(ctx context.Context, wg *sync.WaitGroup, pth string, act NodeAction) {
	defer wg.Done()
	fmt.Printf("Enter %s\n", pth)
	if fi, err := os.Stat(pth); err == nil && act(ctx, pth) && fi.Mode().IsDir() {
		children, err := listChildren(pth)
		if err != nil {
			return
		}
		for _, c := range children {
			wg.Add(1)
			go visit(ctx, wg, filepath.Join(pth, c), act)
		}
	}
	fmt.Printf("Exit %s\n", pth)
}

func listChildren(folder string) ([]string, error) {
	f, err := os.Open(folder)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Readdirnames(0)
}
