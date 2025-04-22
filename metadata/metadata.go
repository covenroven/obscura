package metadata

import "strings"

type Book struct {
	Title        string
	Series       string
	Number       uint32
	Tags         []string
	URLs         []string
	PageCount    uint32
	ParodySeries []string
	Author       Author
	Characters   []string
	Language     string
	Format       string
}

type Author struct {
	Artists    []string
	Circles    []string
	Publishers []string
	Translator []string
}

type FormattedMetadata struct {
	Data     []byte
	FileName string
}

func (b Book) LocalizedTitles() (title, locTitle string) {
	return localized(b.Title)
}

func (b Book) LocalizedSeries() (series, locSeries string) {
	return localized(b.Series)
}

// GetPrimaryAuthor fetches the first author with
// publisher/circle as alternative
func (a Author) GetPrimaryAuthor() (author string) {
	if len(a.Artists) > 0 {
		author = a.Artists[0]

		// discard alternative name
		author, _ = localized(author)
	}

	if author != "" {
		return
	}

	// use circle/publisher as alternative
	if len(a.Publishers) > 0 {
		author = a.Publishers[0]
	}
	if len(a.Circles) > 0 {
		author = a.Circles[0]
	}

	if author == "" {
		author = "[unknown]"
	}

	return
}

// localized will split pipe (|)-separated text into
// original and localized text
func localized(text string) (og, loc string) {
	split := strings.Split(text, "|")
	if len(split) == 0 {
		return
	}

	og = strings.TrimSpace(split[0])
	if len(split) > 1 {
		loc = strings.TrimSpace(split[1])
	}

	return
}
