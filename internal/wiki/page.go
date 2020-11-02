package wiki

const pageElementName = "page"

// Page represents 'page' xml element contents for exported Wikipedia page
// 'page' child elements which not required for parsing operation are omitted
type Page struct {
	Text string `xml:"revision>text" binding:"required,max=255,min=1"`
}
