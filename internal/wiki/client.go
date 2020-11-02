package wiki

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

const wikiPageExportEndpoint = "/wiki/Special:Export/"

// Client fetches and decodes Wikipedia page into Page
type Client struct {
	url string
}

// NewClient initializes Wiki client with given base url
// url parameter must be fully qualified HTTP URL without path or query parameters, example: https://www.somepage.io
func NewClient(url string) *Client {
	return &Client{url: url}
}

// GetPage Fetches and decodes wikipedia page with given page ID
// Will return error when request fails or response is not Wikipedia page
// When err is nil, resp always contains a non-nil Page.Text
func (c *Client) GetPage(pageID string) (*Page, error) {
	url := c.url + wikiPageExportEndpoint + pageID
	logrus.Debugf("fetching exported wiki article from: %s", url)
	response, e := http.Get(url)

	return handleResponse(e, response)
}

func handleResponse(e error, response *http.Response) (*Page, error) {
	if e != nil {
		return nil, fmt.Errorf("fetching exported article encountered error: %w", e)
	}

	defer closeResponse(response)()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wikipage reading failed with error code: %d", response.StatusCode)
	}

	return toResult(response.Body)
}

func closeResponse(response *http.Response) func() {
	return func() {
		if e := response.Body.Close(); e != nil {
			panic(e)
		}
	}
}

func toResult(responseBody io.ReadCloser) (*Page, error) {
	decoder := xml.NewDecoder(responseBody)
	var page Page

	for {
		t, _ := decoder.Token()

		if t == nil {
			break
		}

		se, isStartElement := t.(xml.StartElement)

		if isStartElement && se.Name.Local == pageElementName {
			_ = decoder.DecodeElement(&page, &se)
		}
	}

	if page.Text == "" {
		return nil, errors.New("target page is not wikipedia page")
	}

	return &page, nil
}
