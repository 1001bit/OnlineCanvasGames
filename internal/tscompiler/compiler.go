package tscompiler

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const typeScriptRoot = "web/typescript"

func isPathForbidden(path string) bool {
	forbidden := [...]string{typeScriptRoot, "web/typescript/jquery"}

	for i := range forbidden {
		if path == forbidden[i] {
			return true
		}
	}

	return false
}

func compilePath(path string) {
	if isPathForbidden(path) {
		return
	}

	app := "tsc"
	flag := "-p"

	cmd := exec.Command(app, flag, path)
	out, err := cmd.Output()
	if err != nil {
		log.Println("Error compiling typescript:", string(out))
	}

	log.Println("Compiled typescript from", path)
}

func CompileTypeScript() error {
	// Iterate over directories under typeScriptRoot
	err := filepath.Walk(typeScriptRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			go compilePath(path)
		}

		return nil
	})

	return err
}
