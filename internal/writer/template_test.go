package writer

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/surmus/iso-country-parser/internal"
	"os"
	"testing"
)

func TestNewTemplateResultWriter(t *testing.T) {
	t.Run("should create writer with valid template", func(t *testing.T) {
		writer, err := NewTemplateResultWriter(os.Stdout, "{CODE}, {NAME}")

		assert.NotNil(t, writer)
		assert.Nil(t, err)
	})

	t.Run("should return error when inputting invalid template", func(t *testing.T) {
		writer, err := NewTemplateResultWriter(os.Stdout, "{CODE}, ME}")

		assert.Nil(t, writer)
		assert.NotNil(t, err)
	})
}

func Test_templateResultWriter_Write(t *testing.T) {
	t.Run("should write result", func(t *testing.T) {
		outputBuffer := bytes.NewBuffer([]byte{})

		writer, _ := NewTemplateResultWriter(outputBuffer, "{NAME}, {CODE} ")
		_ = writer.Write([]*internal.Country{
			internal.NewCountry("TEST-1-NAME", "TEST-1-CODE"),
			internal.NewCountry("TEST-2-NAME", "TEST-2-CODE"),
		})

		assert.Equal(t, "TEST-1-NAME, TEST-1-CODE TEST-2-NAME, TEST-2-CODE ", outputBuffer.String())
	})
}
