package formatter

import (
	"context"
	"errors"

	"github.com/covenroven/obscura/metadata"
)

const (
	FormatterComicInfoV2 = "ci2"
)

type Formatter interface {
	// Format converts given metadata into applicable format and return the content
	Format(ctx context.Context, md *metadata.Book) ([]byte, error)

	// FileName returns the file name for the output file including extension
	FileName() string
}

var formatters = make(map[string]Formatter)

func init() {
	formatters[FormatterComicInfoV2] = NewComicInfoV2Formatter()
}

// Format converts given internal metadata to exportable format
func Format(
	ctx context.Context,
	format string,
	md *metadata.Book,
) (*metadata.FormattedMetadata, error) {
	var formatter Formatter
	switch format {
	case FormatterComicInfoV2:
		formatter = formatters[FormatterComicInfoV2]
	default:
		return nil, errors.New("formatter: format not supported")
	}

	formatted, err := formatter.Format(ctx, md)
	if err != nil {
		return nil, err
	}

	return &metadata.FormattedMetadata{
		Data:     formatted,
		FileName: formatter.FileName(),
	}, nil
}
