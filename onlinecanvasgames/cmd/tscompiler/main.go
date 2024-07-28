package main

import (
	"sync"

	"github.com/1001bit/OnlineCanvasGames/pkg/tscompiler"
)

func main() {
	const (
		pagesPath      = "web/typescript/pages"
		gameAssetsPath = "web/typescript/gameassets"
	)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		tscompiler.RecursiveCompileIn(pagesPath)
		wg.Done()
	}()
	go func() {
		tscompiler.RecursiveCompileIn(gameAssetsPath)
		wg.Done()
	}()

	wg.Wait()
}
