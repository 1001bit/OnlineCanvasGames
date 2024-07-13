package tscompiler

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// sourceDir/*.ts -> outputDir/script.js
func compile(sourceDir, outputDir string) error {
	// Ensure the output directory exists
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	// Find all .ts files in the source directory
	tsFiles, err := filepath.Glob(filepath.Join(sourceDir, "*.ts"))
	if err != nil {
		return err
	}

	// Concatenate all TS files into a single temporary file
	tmpFile, err := os.CreateTemp("", "temp-ts-*.ts")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	// sourceDir/*.ts -> tmpFile
	for _, tsFile := range tsFiles {
		content, err := os.ReadFile(tsFile)
		if err != nil {
			return err
		}
		_, err = tmpFile.Write(content)
		if err != nil {
			return err
		}
	}

	// Reset file position to beginning
	tmpFile.Seek(0, 0)

	// tmpFile -> outputDir/script.js
	outputPath := filepath.Join(outputDir, "script.js")
	cmd := exec.Command("tsc", tmpFile.Name(), "-outFile", outputPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}

	log.Printf("Compiled %s -> %s -> %s\n", filepath.Join(sourceDir, "*.ts"), tmpFile.Name(), outputPath)
	return nil
}

func CompileTypeScript() error {
	sourceRoot := "./web/typescript"
	targetRoot := "./web/static"

	// Iterate over directories under sourceRoot
	err := filepath.Walk(sourceRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != sourceRoot {
			pageName := filepath.Base(path)                      // "page"
			targetPageDir := filepath.Join(targetRoot, pageName) // ./web/static/"page"

			// Compile TypeScript files in the current page directory
			err := compile(path, targetPageDir) // ./web/typescript/"page"/*.ts -> ./web/static/"page"/script.js
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
