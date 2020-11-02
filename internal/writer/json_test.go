package writer

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/surmus/iso-country-parser/internal"
	"testing"
)

const expectedJSONResult = `[
	{
		"Name": "TEST-1-NAME",
		"Code": "TEST-1-CODE"
	},
	{
		"Name": "TEST-2-NAME",
		"Code": "TEST-2-CODE"
	}
]
`

func Test_jsonResultWriter_Write(t *testing.T) {
	t.Run("should write result", func(t *testing.T) {
		outputBuffer := bytes.NewBuffer([]byte{})

		writer := NewJSONResultWriter(outputBuffer)
		_ = writer.Write([]*internal.Country{
			internal.NewCountry("TEST-1-NAME", "TEST-1-CODE"),
			internal.NewCountry("TEST-2-NAME", "TEST-2-CODE"),
		})

		assert.Equal(t, expectedJSONResult, outputBuffer.String())
	})
}
