package parser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/covenroven/obscura/metadata"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	hentoidDescriptorFilename = "contentV2.json"

	ehentaiGalleryPrefix  = "https://e-hentai.org/g"
	exhentaiGalleryPrefix = "https://exhentai.org/g"
	nhentaiGalleryPrefix  = "https://nhentai.net/g"
	hitomiGalleryPrefix   = "https://hitomi.la/doujinshi"

	formatOneShot = "One-Shot"
)

type (
	hentoidStdObject struct {
		Name string `json:"name"`
		Type string `json:"type"`
		URL  string `json:"url"`
	}

	hentoidStdObjects []hentoidStdObject

	hentoidDescriptor struct {
		Title         string `json:"title"`
		CoverImageURL string `json:"coverImageUrl"`
		Site          string `json:"site"`
		URL           string `json:"url"` // in /1234567 format
		QtyPages      int64  `json:"qtyPages"`
		Attributes    struct {
			Artists    hentoidStdObjects `json:"ARTIST"`
			Circle     hentoidStdObjects `json:"CIRCLE"`
			Languages  hentoidStdObjects `json:"LANGUAGE"`
			Tags       hentoidStdObjects `json:"TAG"`
			Series     hentoidStdObjects `json:"SERIE"`
			Characters hentoidStdObjects `json:"CHARACTER"`
			Categories hentoidStdObjects `json:"CATEGORY"`
		} `json:"attributes"`
	}
)

type HentoidParser struct{}

var _ Parser = (*HentoidParser)(nil)

func NewHentoidParser() Parser {
	return &HentoidParser{}
}

func (p *HentoidParser) Parse(ctx context.Context, dir *ParsableDir) (*metadata.Book, error) {
	desc, err := p.getDescriptor(dir)
	if err != nil {
		return nil, err
	}

	md := desc.toMetadata()
	return &md, nil
}

func (p *HentoidParser) getDescriptor(dir *ParsableDir) (*hentoidDescriptor, error) {
	log.Println("getting descriptor for", dir.Name)
	var descFile os.DirEntry
	for _, f := range dir.Files {
		if f.Name() == hentoidDescriptorFilename {
			descFile = f
			break
		}
	}

	log.Println("descfile found:", descFile.Name())
	if descFile == nil {
		return nil, errors.New("hentoid: file descriptor not found")
	}

	content, err := os.ReadFile(fmt.Sprintf("%s/%s", dir.Name, descFile.Name()))
	if err != nil {
		log.Println("err open descriptor", err)
		return nil, err
	}

	var desc hentoidDescriptor
	err = json.Unmarshal(content, &desc)
	if err != nil {
		return nil, err
	}

	return &desc, nil
}

func (d hentoidDescriptor) toMetadata() metadata.Book {
	md := metadata.Book{
		Title:        d.Title,
		PageCount:    uint32(d.QtyPages),
		Tags:         d.Attributes.Tags.ExtractName(false),
		URLs:         []string{d.FullURL()},
		ParodySeries: d.Attributes.Series.ExtractName(true),
		Author: metadata.Author{
			Artists: d.Attributes.Artists.ExtractName(true),
			Circles: d.Attributes.Circle.ExtractName(true),
		},
		Characters: d.Attributes.Characters.ExtractName(true),
		Format:     formatOneShot,
	}

	// clean up specific sites
	switch d.Site {
	case "EHENTAI", "EXHENTAI":
		// remove all bracketed text in the title
		re := regexp.MustCompile(`\s*[\(\[\{][^\]\}]*[\)\]\}]\s*`)
		cleanedTitle := strings.TrimSpace(re.ReplaceAllString(md.Title, " "))
		md.Title = cleanedTitle

		// remove tag prefix (male:xyz or female:xyz)
		tags := []string{}
		for _, t := range md.Tags {
			t = strings.ReplaceAll(t, "female:", "")
			t = strings.ReplaceAll(t, "male:", "")
			t = strings.ReplaceAll(t, "mixed:", "")
			t = strings.ReplaceAll(t, "other:", "")
			tags = append(tags, t)
		}
		md.Tags = tags
	}

	// for hentoid series == title
	md.Series = md.Title

	return md
}

// FullURL extracts the book's full URL including the site
func (o hentoidDescriptor) FullURL() (url string) {
	switch o.Site {
	case "EHENTAI":
		url = fmt.Sprintf("%s%s", ehentaiGalleryPrefix, o.URL)
	case "EXHENTAI":
		url = fmt.Sprintf("%s%s", exhentaiGalleryPrefix, o.URL)
	case "NHENTAI":
		url = fmt.Sprintf("%s%s", nhentaiGalleryPrefix, o.URL)
	case "HITOMI":
		url = fmt.Sprintf("%s%s", hitomiGalleryPrefix, o.URL)
	}

	return url
}

// ExtractName extracts the name from each object member
func (o hentoidStdObjects) ExtractName(casing bool) (result []string) {
	for _, obj := range o {
		n := obj.Name
		if casing {
			n = cases.Title(language.English).String(n)
		}
		result = append(result, n)
	}

	return result
}
