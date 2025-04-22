package formatter

import (
	"context"
	"encoding/xml"
	"strings"

	"github.com/covenroven/obscura/metadata"
)

const (
	ComicInfoAgeRatingR18 = "R18+"
)

type ComicInfoV2Formatter struct{}

func NewComicInfoV2Formatter() Formatter {
	return &ComicInfoV2Formatter{}
}

var _ Formatter = (*ComicInfoV2Formatter)(nil)

func (f ComicInfoV2Formatter) Format(ctx context.Context, md *metadata.Book) ([]byte, error) {
	cidata := NewComicInfoV2DataFromMetadata(md)

	result, err := xml.Marshal(cidata)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (f ComicInfoV2Formatter) FileName() string {
	return "ComicInfo.xml"
}

type comicInfoV2Data struct {
	XMLName xml.Name `xml:"ComicInfo"`

	Title           string `xml:"Title"`
	Series          string `xml:"Series,omitempty"`
	LocalizedSeries string `xml:"LocalizedSeries,omitempty"` // non-standard
	SeriesGroup     string `xml:"SeriesGroup,omitempty"`
	Count           uint32 `xml:"Count,omitempty"`
	Volume          uint32 `xml:"Volume,omitempty"`
	Tags            string `xml:"Tags"`
	Web             string `xml:"Web"`
	PageCount       uint32 `xml:"PageCount"`
	AgeRating       string `xml:"AgeRating"`
	Characters      string `xml:"Characters"`
	Format          string `xml:"Format"`

	Writer      string `xml:"Writer"`
	CoverArtist string `xml:"CoverArtist"`
	Publisher   string `xml:"Publisher,omitempty"`
	Translator  string `xml:"Translator,omitempty"`
}

func NewComicInfoV2DataFromMetadata(md *metadata.Book) *comicInfoV2Data {
	ci := new(comicInfoV2Data)
	ci.Title = md.Title
	ci.Series, ci.LocalizedSeries = md.LocalizedSeries()
	ci.SeriesGroup = strings.Join(md.ParodySeries, ",")
	ci.Volume = md.Number
	ci.Tags = strings.Join(md.Tags, ",")
	ci.Web = strings.Join(md.URLs, ",")
	ci.PageCount = md.PageCount
	ci.AgeRating = ComicInfoAgeRatingR18
	ci.Writer = strings.Join(md.Author.Artists, ",")
	ci.Publisher = strings.Join(md.Author.Publishers, ",")
	if len(md.Author.Circles) > 0 {
		ci.Publisher = strings.Join(md.Author.Circles, ",")
	}
	ci.Translator = strings.Join(md.Author.Translator, ",")
	ci.Characters = strings.Join(md.Characters, ",")
	ci.Format = md.Format

	// mark as series ended
	ci.Count = 1

	return ci
}
