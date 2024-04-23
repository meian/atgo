package roundtrippers

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httputil"
)

type LoggingRoundTripper struct {
	transport http.RoundTripper
	logger    slog.Logger
}

var _ http.RoundTripper = &LoggingRoundTripper{}

func NewLoggingRoundTripper(transport http.RoundTripper, logger slog.Logger) http.RoundTripper {
	if !logger.Enabled(context.Background(), slog.LevelDebug) {
		return transport
	}
	return &LoggingRoundTripper{
		transport: transport,
		logger:    logger,
	}
}

func (rt *LoggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	requestDump, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		rt.logger.Error(err.Error())
		return nil, err
	}
	rt.logger.Debug(string(requestDump))

	response, err := rt.transport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	responseDump, err := httputil.DumpResponse(response, true)
	if err != nil {
		rt.logger.Error(err.Error())
		return nil, err
	}
	rt.logger.Debug(string(responseDump))

	return response, nil
}
