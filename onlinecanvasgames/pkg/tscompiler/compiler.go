package tscompiler

import (
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// Compile single tsconfig inside path
func executeTscIn(path string) {
	tsconfigPath := filepath.Join(path, "tsconfig.json")
	if _, err := os.Stat(tsconfigPath); err != nil {
		return
	}

	app := "tsc"
	flag := "-p"

	cmd := exec.Command(app, flag, tsconfigPath)
	out, err := cmd.Output()
	if err != nil {
		log.Println("Error compiling typescript:", err, string(out))
	}

	log.Println("Compiled typescript from", path)
}

// Compile all the tsconfigs inside all the directories inside path
func RecursiveCompileIn(rootPath string) {
	var wg sync.WaitGroup

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() || path == rootPath {
			return nil
		}

		wg.Add(1)
		go func() {
			executeTscIn(path)
			wg.Done()
		}()

		return nil
	})

	if err != nil {
		log.Println("WalkDir error:", err)
	}

	wg.Wait()
}
