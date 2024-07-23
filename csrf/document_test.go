package csrf_test

import (
	"embed"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/meian/atgo/csrf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/document
	testdata embed.FS
)

func TestFromDocument(t *testing.T) {
	htmlm := testHTMLMap(t)

	tests := []struct {
		name string
		want string
	}{
		{name: "valid", want: "test_csrf_token"},
		{name: "no-head", want: ""},
		{name: "no-script", want: ""},
		{name: "no-token", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlm[tt.name+".html"]))
			require.NoError(t, err)
			assert.Equal(t, tt.want, csrf.FromDocument(doc))
		})
	}
}

func testHTMLMap(t *testing.T) map[string]string {
	t.Helper()
	m := make(map[string]string)
	es, err := testdata.ReadDir("testdata/document")
	require.NoError(t, err)
	for _, e := range es {
		require.False(t, e.IsDir())
		b, err := testdata.ReadFile("testdata/document/" + e.Name())
		require.NoError(t, err)
		m[e.Name()] = string(b)
	}
	return m
}
