package parser

import (
	"github.com/sirupsen/logrus"
	"github.com/surmus/iso-country-parser/internal"
	"github.com/surmus/iso-country-parser/internal/wiki"
	"regexp"
)

// matches contents of [[CONTENTS]]
const interWikiLinkDeclarationPattern = `\[\[\s*((.|\n|\r)*?)]]`

// WikiPageParser converts data from Wiki page into countries list
type WikiPageParser struct {
	page *wiki.Page

	interWikiLinkMatcher *regexp.Regexp
	countryCodeMatcher   *regexp.Regexp
}

// NewWikiPageParser creates new parser for given Wiki page
func NewWikiPageParser(page *wiki.Page) *WikiPageParser {
	return &WikiPageParser{
		interWikiLinkMatcher: regexp.MustCompile(interWikiLinkDeclarationPattern),
		countryCodeMatcher:   regexp.MustCompile(isoAlpha3CountryCodePattern),
		page:                 page,
	}
}

// Parse parses Wiki page contents into list of countries, if no countries are found from page empty slice is returned
func (p *WikiPageParser) Parse() []*internal.Country {
	logrus.Info("begin parsing wiki page content")
	links := p.interWikiLinkMatcher.FindAllStringSubmatch(p.page.Text, -1)
	logrus.Debugf("found %d interwiki links", len(links))
	linkParser := newInterWikiLinkParser(links)

	return linkParser.parse()
}
