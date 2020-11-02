package writer

import (
	"fmt"
	"github.com/surmus/iso-country-parser/internal"
	"io"
	"strings"
	"text/template"
)

const codePlaceholder = "{CODE}"
const namePlaceHolder = "{NAME}"

type templateResultWriter struct {
	template    string
	destination io.Writer
}

// NewTemplateResultWriter creates ResultWriter which formats countries into text according to supplied template
// template argument describes how internal.Country is represented in textual form
// example template: '{CODE}, {NAME}' will produce 'USA, The United States of America'
// will return error if template does not contain any either placeholders({CODE} or {NAME})
func NewTemplateResultWriter(destination io.Writer, template string) (ResultWriter, error) {
	if e := validateTemplate(template); e != nil {
		return nil, e
	}

	return &templateResultWriter{template: buildTemplate(template), destination: destination}, nil
}

func (w *templateResultWriter) Write(countries []*internal.Country) error {
	tpl := template.Must(template.New("sql").Parse(w.template))

	return tpl.Execute(w.destination, countries)
}

func validateTemplate(template string) error {
	isValid := strings.Contains(template, codePlaceholder) || strings.Contains(template, namePlaceHolder)

	if isValid {
		return nil
	}

	return fmt.Errorf("invalid result template")
}

func buildTemplate(userProviderTemplate string) string {
	templateText := strings.ReplaceAll(userProviderTemplate, codePlaceholder, "{{ .Code }}")
	templateText = strings.ReplaceAll(templateText, namePlaceHolder, "{{ .Name }}")

	return "{{ range . }}" + templateText + "{{ end }}"
}
