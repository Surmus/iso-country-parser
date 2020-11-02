package writer

import (
	"github.com/surmus/iso-country-parser/internal"
)

// ResultWriter converts countries into textual form
type ResultWriter interface {
	Write(countries []*internal.Country) error
}
