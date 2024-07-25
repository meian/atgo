package crawler_test

import (
	"embed"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type htmlMap map[string]string

func (m htmlMap) Get(key string) string {
	if key == "no-html" {
		return "no html"
	}
	return m[key]
}

var (
	//go:embed testdata
	testdata embed.FS
	htmlmap  = map[string]htmlMap{}
)

func testHTMLMap(t *testing.T, target string) htmlMap {
	if m, ok := htmlmap[target]; ok {
		return m
	}
	t.Helper()
	m := make(map[string]string)
	dir := path.Join("testdata", target)
	es, err := testdata.ReadDir(dir)
	require.NoError(t, err)
	for _, e := range es {
		require.False(t, e.IsDir())
		b, err := testdata.ReadFile(path.Join(dir, e.Name()))
		require.NoError(t, err)
		m[e.Name()] = string(b)
	}
	htmlmap[target] = m
	return m
}

type mockRequestRoundTripper struct {
	request *http.Request
}

func (m *mockRequestRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.request = req
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("OK")),
	}, nil
}

type captureFunc func() (method, path string, query, body url.Values)

func (m *mockRequestRoundTripper) lastCaputure() (string, string, url.Values, url.Values) {
	query := m.request.URL.Query()
	body := url.Values{}
	if m.request.Body != nil {
		b, _ := io.ReadAll(m.request.Body)
		if bt, err := url.ParseQuery(string(b)); err == nil {
			body = bt
		} else {
			panic(errors.Wrapf(err, "cannot parse request body: %s", string(b)))
		}
	}
	return m.request.Method, m.request.URL.Path, query, body
}

func mockRequestClient() (*http.Client, captureFunc) {
	m := &mockRequestRoundTripper{}
	c := &http.Client{
		Transport: m,
	}
	return c, m.lastCaputure
}

type mockResponseRoundTripper struct {
	status  int
	body    string
	timeout bool
}

func (m *mockResponseRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var wait <-chan time.Time
	if m.timeout {
		wait = time.After(1 * time.Second)
	} else {
		wait = time.After(0 * time.Second)
	}
	select {
	case <-wait:
		return &http.Response{
			StatusCode: m.status,
			Body:       io.NopCloser(strings.NewReader(m.body)),
		}, nil
	case <-req.Context().Done():
		return nil, errors.New("request canceled for timeout")
	}
}

func mockResponseClient(status int, body string, timeout bool) *http.Client {
	return &http.Client{
		Transport: &mockResponseRoundTripper{status: status, body: body, timeout: timeout},
	}
}
