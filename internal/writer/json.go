package writer

import (
	"encoding/json"
	"github.com/surmus/iso-country-parser/internal"
	"io"
)

type jsonResultWriter struct {
	destination io.Writer
}

// NewJSONResultWriter creates ResultWriter which formats countries JSON list
func NewJSONResultWriter(destination io.Writer) ResultWriter {
	return &jsonResultWriter{destination: destination}
}

func (w jsonResultWriter) Write(countries []*internal.Country) error {
	encoder := json.NewEncoder(w.destination)
	encoder.SetIndent("", "\t")

	return encoder.Encode(countries)
}
