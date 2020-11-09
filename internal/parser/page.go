package parser

import (
	"github.com/dlclark/regexp2"
	"github.com/sirupsen/logrus"
	"github.com/surmus/iso-country-parser/internal"
	"github.com/surmus/iso-country-parser/internal/wiki"
)

// matches strings which start with 'id=SOMEID' and end with '[[ISO 3166-2:XX]]' 'XX' being alpha-2 country code
const countriesTableRowPattern = `id=\s*((?!id=).|\n|\r)*?\[\[ISO 3166-2:[A-Z]{2}]]`

// WikiPageParser converts data from Wiki page into countries list
type WikiPageParser struct {
	page                     *wiki.Page
	countryRowParser         *countriesTableRowParser
	countriesTableRowMatcher *regexp2.Regexp
}

// NewWikiPageParser creates new parser for given Wiki page
func NewWikiPageParser(page *wiki.Page) *WikiPageParser {
	return &WikiPageParser{
		countriesTableRowMatcher: regexp2.MustCompile(countriesTableRowPattern, regexp2.ECMAScript),
		page:                     page,
		countryRowParser:         newCountriesTableRowParser(),
	}
}

// Parse parses Wiki page contents into list of countries, if no countries are found from page empty slice is returned
func (p *WikiPageParser) Parse() []*internal.Country {
	logrus.Info("begin parsing wiki page content")
	countryRowMatch, _ := p.countriesTableRowMatcher.FindStringMatch(p.page.Text)

	results := p.appendCountryMatch(countryRowMatch, []*internal.Country{})
	logrus.Infof("found %d countries from wiki page", len(results))
	return results
}

func (p *WikiPageParser) appendCountryMatch(match *regexp2.Match, results []*internal.Country) []*internal.Country {
	if match == nil {
		return results
	}

	if country := p.countryRowParser.parse(match.String()); country != nil {
		results = append(results, country)
	}

	nextMatch, _ := p.countriesTableRowMatcher.FindNextMatch(match)
	return p.appendCountryMatch(nextMatch, results)
}
