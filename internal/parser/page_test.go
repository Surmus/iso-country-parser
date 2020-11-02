package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/surmus/iso-country-parser/internal"
	"github.com/surmus/iso-country-parser/internal/wiki"
	"testing"
)

const testWikiPageContent = `|id=VENE|{{flagicon|Venezuela|state}}&nbsp;[[Venezuela|Venezuela (Bolivarian Republic of)]] |[[Venezuela|The Bolivarian Republic of Venezuela]] |align=center|UN member state |align=center|[[ISO 3166-1 alpha-2#VE|{{mono|VE}}]] |align=center|[[ISO 3166-1 alpha-3#VEN|{{mono|VEN}}]] |align=center|[[ISO 3166-1 numeric#862|{{mono|862}}]] |align=center|[[ISO 3166-2:VE]] |align=center|[[.ve]] |-style="background-color:white;"|id=AZER|{{flagicon|Azerbaijan}}&nbsp;[[Azerbaijan]] |[[Azerbaijan|The Republic of Azerbaijan]] |align=center|UN member state |align=center|[[ISO 3166-1 alpha-2#AZ|{{mono|AZ}}]] |align=center|[[ISO 3166-1 alpha-3#AZE|{{mono|AZE}}]] |align=center|[[ISO 3166-1 numeric#031|{{mono|031}}]] |align=center|[[ISO 3166-2:AZ]] |align=center|[[.az]] |-style="background-color:white;`
const invalidWikiPageContent = `[[ISO 3166-1 alpha-3#AZE|{{mono|AZE}}]]`

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

	t.Run("should not extract partial matches", func(t *testing.T) {
		parser := NewWikiPageParser(&wiki.Page{Text: invalidWikiPageContent})

		assert.Empty(t, parser.Parse())
	})

}
