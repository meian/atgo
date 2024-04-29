package http

import (
	"context"
	"net/http"
)

type key string

const (
	clientKey   key = "client"
	skipWaitKey key = "skipWait"
)

func ClientFromContext(ctx context.Context) *http.Client {
	c, ok := ctx.Value(clientKey).(*http.Client)
	if !ok {
		return http.DefaultClient
	}
	return c
}

func ContextWithClient(ctx context.Context, client *http.Client) context.Context {
	return context.WithValue(ctx, clientKey, client)
}

func IsSkipWait(ctx context.Context) bool {
	s, ok := ctx.Value(skipWaitKey).(bool)
	if !ok {
		return false
	}
	return s
}

func ContextWithSkipWait(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipWaitKey, true)
}
