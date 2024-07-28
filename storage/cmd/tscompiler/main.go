package main

import (
	"sync"

	"github.com/neinBit/ocg-storage-service/pkg/tscompiler"
)

func main() {
	const (
		pagesPath      = "typescript/pages"
		gameAssetsPath = "typescript/gameassets"
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
