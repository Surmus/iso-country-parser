package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/surmus/iso-country-parser/internal"
	"github.com/surmus/iso-country-parser/internal/wiki"
	"testing"
)

const testWikiPageContent = `|id=VENE|{{flagicon|Venezuela|state}}&nbsp;[[Venezuela|Venezuela (Bolivarian Republic of)]] |[[Venezuela|The Bolivarian Republic of Venezuela]] |align=center|UN member state |align=center|[[ISO 3166-1 alpha-2#VE|{{mono|VE}}]] |align=center|[[ISO 3166-1 alpha-3#VEN|{{mono|VEN}}]] |align=center|[[ISO 3166-1 numeric#862|{{mono|862}}]] |align=center|[[ISO 3166-2:VE]] |align=center|[[.ve]] |-style="background-color:white;"|id=AZER|{{flagicon|Azerbaijan}}&nbsp;[[Azerbaijan]] |[[Azerbaijan|The Republic of Azerbaijan]] |align=center|UN member state |align=center|[[ISO 3166-1 alpha-2#AZ|{{mono|AZ}}]] |align=center|[[ISO 3166-1 alpha-3#AZE|{{mono|AZE}}]] |align=center|[[ISO 3166-1 numeric#031|{{mono|031}}]] |align=center|[[ISO 3166-2:AZ]] |align=center|[[.az]] |-style="background-color:white;`

const wikiPageCountryRowWithManyNameParts = `id=ANTA|{{flagicon|Antarctica}}&nbsp;
[[Antarctica]] 
{{efn|The [[ISO 3166]] country name [[Antarctica]] comprises the [[Antarctica|continent of Antarctica]] and all land and ice shelves south of the [[60th parallel south]].}} 
|All land and ice shelves south of the [[60th parallel south]] 
|align=center|[[Antarctic Treaty System|Antarctic Treaty]] 
|align=center|[[ISO 3166-1 alpha-2#AQ|{{mono|AQ}}]] 
|align=center|[[ISO 3166-1 alpha-3#ATA|{{mono|ATA}}]] 
|align=center|[[ISO 3166-1 numeric#010|{{mono|010}}]] 
|align=center|[[ISO 3166-2:AQ]] |align=center|[[.aq]] 
|-style="background-color:white;"`

func TestWikiPageParser_Parse(t *testing.T) {
	t.Run("should parse wiki page content into results", func(t *testing.T) {
		parser := NewWikiPageParser(&wiki.Page{Text: testWikiPageContent})

		assert.Equal(
			t,
			[]*internal.Country{
				internal.NewCountry("Venezuela", "VEN"),
				internal.NewCountry("Azerbaijan", "AZE"),
			},
			parser.Parse(),
		)
	})

	t.Run("should parse wiki page row with many name parts into results", func(t *testing.T) {
		parser := NewWikiPageParser(&wiki.Page{Text: wikiPageCountryRowWithManyNameParts})

		assert.Equal(
			t,
			[]*internal.Country{
				internal.NewCountry("Antarctica", "ATA"),
			},
			parser.Parse(),
		)
	})

	t.Run("should ignore partial matches", func(t *testing.T) {
		parser := NewWikiPageParser(&wiki.Page{Text: `id=ANTA|{{flagicon|Antarctica}}&nbsp|align=center|[[ISO 3166-2:TE]]`})

		assert.Empty(t, parser.Parse())
	})

	t.Run("should ignore table row without ISO alpha 3 code", func(t *testing.T) {
		parser := NewWikiPageParser(&wiki.Page{Text: `id=ANTA|{{flagicon|Antarctica}}&nbsp[[Antarctica]] |align=center|[[ISO 3166-2:AQ]] |align=center|[[.aq]]`})

		assert.Empty(t, parser.Parse())
	})
}
