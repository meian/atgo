package http

import (
	"context"
	"net/http"
)

type key string

const clientKey key = "client"

func ClientFromContext(ctx context.Context) *http.Client {
	v := ctx.Value(clientKey)
	c, ok := v.(*http.Client)
	if !ok {
		return http.DefaultClient
	}
	return c
}

func ContextWithClient(ctx context.Context, client *http.Client) context.Context {
	return context.WithValue(ctx, clientKey, client)
}
