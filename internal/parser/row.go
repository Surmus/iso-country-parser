package parser

import (
	"github.com/sirupsen/logrus"
	"github.com/surmus/iso-country-parser/internal"
	"regexp"
	"strings"
)

// matches EST from `ISO 3166-1 alpha-3#EST|{{mono|EST}}`
const isoAlpha3CountryCodePattern = `(?:ISO 3166-1 alpha-3#)([A-Z0-9]{3})`

// matches contents of [[CONTENTS]]
const interWikiLinkDeclarationPattern = `\[\[\s*((.|\n|\r)*?)]]`
const wikiTableCellSeparator = "|"

const countryNameLinkPosition = 0
const firstMatchGroup = 1

type countriesTableRowParser struct {
	countryCodeMatcher              *regexp.Regexp
	interWikiLinkDeclarationMatcher *regexp.Regexp
}

func newCountriesTableRowParser() *countriesTableRowParser {
	return &countriesTableRowParser{
		countryCodeMatcher:              regexp.MustCompile(isoAlpha3CountryCodePattern),
		interWikiLinkDeclarationMatcher: regexp.MustCompile(interWikiLinkDeclarationPattern),
	}
}

// parse validates and extracts country information from given countries table row
// nil is returned when row does not contain required data
func (p *countriesTableRowParser) parse(countriesTableRow string) *internal.Country {
	rowInterWikiLinkContents := p.extractLinks(countriesTableRow)

	if len(rowInterWikiLinkContents) < 2 {
		logrus.Warnf("found invalid country table row: %s", countriesTableRow)
		return nil
	}

	countryName := p.extractCountryName(rowInterWikiLinkContents)

	if countryCode := p.extractCountryCode(rowInterWikiLinkContents); countryCode != "" {
		return internal.NewCountry(countryName, countryCode)
	}
	logrus.Warnf("found row without country code: %s", countriesTableRow)
	return nil
}

func (p *countriesTableRowParser) extractLinks(countriesTableRow string) []string {
	var links []string

	for _, linkMatch := range p.interWikiLinkDeclarationMatcher.FindAllStringSubmatch(countriesTableRow, -1) {
		links = append(links, linkMatch[firstMatchGroup])
	}

	return links
}

func (p *countriesTableRowParser) extractCountryCode(linksContents []string) string {
	for _, linkContents := range linksContents {
		if countryCodeMatch := p.countryCodeMatcher.FindStringSubmatch(linkContents); countryCodeMatch != nil {
			return countryCodeMatch[firstMatchGroup]
		}
	}

	return ""
}

func (p *countriesTableRowParser) extractCountryName(linksContents []string) string {
	nameParts := strings.Split(linksContents[countryNameLinkPosition], wikiTableCellSeparator)

	return strings.TrimSpace(nameParts[0])
}
