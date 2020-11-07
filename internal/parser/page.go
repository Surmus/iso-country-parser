package parser

import (
	"github.com/sirupsen/logrus"
	"github.com/surmus/iso-country-parser/internal"
	"github.com/surmus/iso-country-parser/internal/wiki"
	"regexp"
)

// matches strings which start with 'id=SOMEID' and end with '[[ISO 3166-2:XX]]' 'XX' being alpha-2 country code
const countriesTableRowPattern = `id=\s*(.|\n|\r)*?\[\[ISO 3166-2:[A-Z]{2}]]`

// WikiPageParser converts data from Wiki page into countries list
type WikiPageParser struct {
	page                     *wiki.Page
	countryRowParser         *countriesTableRowParser
	countriesTableRowMatcher *regexp.Regexp
}

// NewWikiPageParser creates new parser for given Wiki page
func NewWikiPageParser(page *wiki.Page) *WikiPageParser {
	return &WikiPageParser{
		countriesTableRowMatcher: regexp.MustCompile(countriesTableRowPattern),
		page:                     page,
		countryRowParser:         newCountriesTableRowParser(),
	}
}

// Parse parses Wiki page contents into list of countries, if no countries are found from page empty slice is returned
func (p *WikiPageParser) Parse() []*internal.Country {
	var result []*internal.Country
	logrus.Info("begin parsing wiki page content")
	countryRows := p.countriesTableRowMatcher.FindAllString(p.page.Text, -1)

	for _, row := range countryRows {
		if country := p.countryRowParser.parse(row); country != nil {
			result = append(result, country)
		}
	}
	logrus.Infof("found %d countries from wiki page", len(result))
	return result
}
