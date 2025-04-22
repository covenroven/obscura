package parser

import (
	"context"
	"errors"
	"os"

	"github.com/covenroven/obscura/metadata"
)

const (
	ParserHentoid = "hentoid"
)

type Parser interface {
	// Parse read the directory and generate new book metadata
	Parse(ctx context.Context, dir *ParsableDir) (*metadata.Book, error)
}

var parsers = make(map[string]Parser)

func init() {
	parsers[ParserHentoid] = NewHentoidParser()
}

// ParsableDir defines a book directory containing its path
// and the containing files
type ParsableDir struct {
	// Name defines the name of the directory
	Name string
	// Files contains every files in the directory
	Files []os.DirEntry
}

func NewParsableDir(name string, files []os.DirEntry) *ParsableDir {
	return &ParsableDir{
		Name:  name,
		Files: files,
	}
}

func Parse(ctx context.Context, source string, dir *ParsableDir) (*metadata.Book, error) {
	switch source {
	case ParserHentoid:
		return parsers[ParserHentoid].Parse(ctx, dir)
	default:
		return nil, errors.New("parser: source not supported")
	}

	return nil, errors.New("parser: source not supported")
}
