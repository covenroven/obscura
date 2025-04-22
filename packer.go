package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

const (
	outputExt = "cbz"
)

var (
	ignoreFiles = []string{
		"thumb.webp",
		"thumb.jpg",
		"thumb.png",
		"thumb.jpeg",
	}
)

func packDir(ctx context.Context, dirName, outputDir, filename string) error {
	if err := mkdir(outputDir); err != nil {
		return err
	}

	return compress(dirName, outputDir, filename)
}

func compress(dirName, outputDir, filename string) error {
	// sp := strings.Split(dirName, "/")
	// filename := sp[len(sp)-1]

	zipfile, err := os.Create(
		fmt.Sprintf("%s/%s.%s", outputDir, filename, outputExt),
	)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	zw := zip.NewWriter(zipfile)
	defer zw.Close()

	if err := filepath.WalkDir(
		dirName,
		func(path string, d fs.DirEntry, err error) error {
			log.Println("crawling", path)
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			if shouldIgnoreFile(d.Name()) {
				log.Println("ignoring", d.Name())
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			targetFile, err := zw.Create(d.Name())
			if err != nil {
				return err
			}

			_, err = io.Copy(targetFile, file)
			if err != nil {
				return err
			}

			return nil
		},
	); err != nil {
		return err
	}

	return nil
}

func writeFile(content []byte, filename string) error {
	return os.WriteFile(filename, content, 0666)
}

func mkdir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func shouldIgnoreFile(filename string) bool {
	for _, ignoref := range ignoreFiles {
		if ignoref == filename {
			return true
		}
	}

	return false
}
