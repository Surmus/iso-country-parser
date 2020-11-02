package internal

// Country is extracted during parse operation from WIKI page
type Country struct {
	Name string
	Code string
}

// NewCountry initializes Country type from officially recognized country name and ISO 3166 alpha-3 country code
func NewCountry(name string, code string) *Country {
	return &Country{Name: name, Code: code}
}
