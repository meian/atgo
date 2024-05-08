package roundtrippers

import (
	gohttp "net/http"
	"time"

	"github.com/meian/atgo/http"
)

type RateLimitRoundTripper struct {
	transport gohttp.RoundTripper
	interval  time.Duration
	lastTime  time.Time
}

var _ gohttp.RoundTripper = &RateLimitRoundTripper{}

func NewRateLimitRoundTripper(transport gohttp.RoundTripper, interval time.Duration) gohttp.RoundTripper {
	if interval <= 0 {
		return transport
	}
	return &RateLimitRoundTripper{
		transport: transport,
		interval:  interval,
	}
}

func (rt *RateLimitRoundTripper) RoundTrip(req *gohttp.Request) (*gohttp.Response, error) {
	skipWait := http.IsSkipWait(req.Context())
	if !skipWait {
		wait := time.Until(rt.lastTime.Add(rt.interval))
		if wait > 0 {
			time.Sleep(wait)
		}
	}
	next := time.Now()
	resp, err := rt.transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	// レスポンスがリダイレクトの場合は次回を待機しない
	// => リダイレクトでなければ正常も異常も次回は待機対象
	if resp.StatusCode != gohttp.StatusFound && !skipWait {
		rt.lastTime = next
	}
	return resp, nil
}
