package roundtrippers

import (
	"net/http"
	"time"
)

type RateLimitRoundTripper struct {
	transport http.RoundTripper
	lastTime  time.Time
	interval  time.Duration
}

var _ http.RoundTripper = &RateLimitRoundTripper{}

func NewRateLimitRoundTripper(transport http.RoundTripper, interval time.Duration) http.RoundTripper {
	if interval <= 0 {
		return transport
	}
	return &RateLimitRoundTripper{
		transport: transport,
		interval:  interval,
	}
}

func (rt *RateLimitRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	now := time.Now()
	wait := rt.lastTime.Add(rt.interval).Sub(now)
	if wait > 0 {
		time.Sleep(wait)
	}
	rt.lastTime = time.Now()
	return rt.transport.RoundTrip(req)
}
