package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/covenroven/obscura/formatter"
	"github.com/covenroven/obscura/parser"
	"golang.org/x/sync/errgroup"
)

var command = "repack"
var format = "ci2"
var source = "hentoid"

var sourcedir, targetdir string

func init() {
	flag.StringVar(&sourcedir, "source", "", "the source directory to act upon")
	flag.StringVar(&targetdir, "out", "", "the output directory for repacked files")
}

func main() {
	flag.Parse()

	if sourcedir == "" || targetdir == "" {
		log.Fatal("source and out argument not supplied")
	}

	if err := Repack(); err != nil {
		log.Fatal(err)
	}
}

func Repack() error {
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return processDir(ctx, sourcedir)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func processDir(ctx context.Context, dir string) error {
	// scan directory
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var hasDir bool
	for _, f := range files {
		if hasDir = f.IsDir(); hasDir {
			break
		}
	}

	if !hasDir {
		// pack current dir
		return identifyAndPack(ctx, dir, files)
	}

	// if contains directory, treat as another directory
	// and process it in another goroutine
	g, ctx := errgroup.WithContext(ctx)

	for _, f := range files {
		g.Go(func() error {
			// no-op for files
			if !f.IsDir() {
				return nil
			}

			return processDir(ctx, fmt.Sprintf("%s/%s", dir, f.Name()))
		})
	}

	if err := g.Wait(); err != nil {
		log.Println(err)
	}

	return err
}

// identifyAndPack optionally creates book metadata and
// re-pack the entries to new directory
func identifyAndPack(ctx context.Context, dirName string, files []os.DirEntry) error {
	dir := parser.NewParsableDir(dirName, files)
	md, err := parser.Parse(ctx, source, dir)
	if err != nil {
		log.Println("skipping", dirName)
		return err
	}

	// ignore invalid metadata
	if md == nil {
		return nil
	}

	formattedMD, err := formatter.Format(ctx, format, md)
	if err != nil {
		return err
	}

	// write the exportable metadata to the current dir
	if err := writeFile(
		formattedMD.Data,
		fmt.Sprintf("%s/%s", dirName, formattedMD.FileName),
	); err != nil {
		return err
	}

	// target base dir: <root>/<artist>/
	targetBaseDir := fmt.Sprintf("%s/%s", targetdir, md.Author.GetPrimaryAuthor())

	filename, _ := md.LocalizedTitles()
	return packDir(ctx, dirName, targetBaseDir, filename)
}
