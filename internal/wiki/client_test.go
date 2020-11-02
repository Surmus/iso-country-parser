package wiki

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

const invalidResponse = `<?xml version="1.0" encoding="UTF-8"?>
<note>
  <to>Tove</to>
  <from>Jani</from>
  <heading>Reminder</heading>
  <body>Don't forget me this weekend!</body>
</note>`

func TestGetPage(t *testing.T) {
	t.Run("should fetch and decode wikipedia page", func(t *testing.T) {
		testResponseFile, _ := os.Open("_test_wikipage_response.xml")
		testResponseFileContents, _ := ioutil.ReadAll(testResponseFile)

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			_, _ = res.Write(testResponseFileContents)
		}))
		defer func() { testServer.Close() }()

		page, err := NewClient(testServer.URL).GetPage("TEST")

		assert.NotNil(t, page)
		assert.Nil(t, err)
		assert.NotEmpty(t, page.Text)
	})

	t.Run("should return error when retrieved contents is not ISO countries page", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			_, _ = res.Write([]byte(invalidResponse))
		}))
		defer func() { testServer.Close() }()

		page, err := NewClient(testServer.URL).GetPage("TEST")

		assert.Nil(t, page)
		assert.EqualError(t, err, "target page is not wikipedia page")
	})

	t.Run("should return error when request fails", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(500)
			_, _ = res.Write([]byte("ERROR"))
		}))
		defer func() { testServer.Close() }()

		page, err := NewClient(testServer.URL).GetPage("TEST")

		assert.Nil(t, page)
		assert.EqualError(t, err, "wikipage reading failed with error code: 500")
	})

	t.Run("should return error when request fails", func(t *testing.T) {
		page, err := NewClient("INVALID-URL").GetPage("TEST")

		assert.Nil(t, page)
		urlErr := &url.Error{}
		assert.True(t, errors.As(err, &urlErr))
	})
}
