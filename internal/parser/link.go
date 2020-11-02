package parser

import (
	"github.com/sirupsen/logrus"
	"github.com/surmus/iso-country-parser/internal"
	"regexp"
	"strings"
)

// matches EST from `ISO 3166-1 alpha-3#EST|{{mono|EST}}`
const isoAlpha3CountryCodePattern = `(?:ISO 3166-1 alpha-3#)([A-Z0-9]{3})`
const wikiTableCellSeparator = "|"
const countryNameRelativePosition = 3

const firstMatchGroup = 1

type interWikiLinkParser struct {
	linkMatches [][]string

	countryCodeMatcher *regexp.Regexp
}

func newInterWikiLinkParser(linkMatches [][]string) *interWikiLinkParser {
	return &interWikiLinkParser{
		linkMatches:        linkMatches,
		countryCodeMatcher: regexp.MustCompile(isoAlpha3CountryCodePattern),
	}
}

func (p *interWikiLinkParser) parse() []*internal.Country {
	var results []*internal.Country

	for i, linkMatch := range p.linkMatches {
		if country := p.parseLink(linkMatch, i); country != nil {
			results = append(results, country)
		} else {
			logrus.Tracef("skipping non iso code link match: %s", linkMatch)
		}
	}
	logrus.Infof("parsed %d interwiki links into %d iso countries", len(p.linkMatches), len(results))
	return results
}

func (p *interWikiLinkParser) parseLink(linkMatch []string, linkPosition int) *internal.Country {
	countryCode := p.extractCountryCode(linkMatch[firstMatchGroup])

	if countryCode == "" {
		logrus.Tracef("skipping non iso code link match: %s", linkMatch)
		return nil
	}

	if countryName := p.extractCountryName(linkPosition); countryName != "" {
		return internal.NewCountry(countryName, countryCode)
	}
	logrus.Debugf("encountered country code field without matching name field: %s", linkMatch)
	return nil
}

func (p *interWikiLinkParser) extractCountryCode(input string) string {
	countryCodeMatchGroups := p.countryCodeMatcher.FindStringSubmatch(input)

	if len(countryCodeMatchGroups) == 0 {
		return ""
	}

	return countryCodeMatchGroups[firstMatchGroup]
}

func (p *interWikiLinkParser) extractCountryName(countryCodePosition int) string {
	namePosition := countryCodePosition - countryNameRelativePosition

	if namePosition < 0 || len(p.linkMatches) < namePosition+1 {
		return ""
	}

	nameParts := strings.Split(p.linkMatches[namePosition][firstMatchGroup], wikiTableCellSeparator)

	return strings.TrimSpace(nameParts[0])
}
