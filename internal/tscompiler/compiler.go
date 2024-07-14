package tscompiler

import (
	"errors"
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

func CompileTypeScript() error {
	// Iterate over directories under sourceRoot
	err := filepath.Walk(typeScriptRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() || isPathForbidden(path) {
			return nil
		}

		app := "tsc"
		flag := "-p"

		cmd := exec.Command(app, flag, path)
		out, err := cmd.Output()
		if err != nil {
			return errors.New(string(out))
		}

		log.Println("compiled typescript from", path)

		return nil
	})

	return err
}
