package io

import (
	"context"
	"io"
)

type key string

const (
	outKey key = "out"
	errKey key = "err"
)

func OutFromContext(ctx context.Context) io.Writer {
	return ctx.Value(outKey).(io.Writer)
}

func OutWithContext(ctx context.Context, out io.Writer) context.Context {
	return context.WithValue(ctx, outKey, out)
}

func ErrFromContext(ctx context.Context) io.Writer {
	return ctx.Value(errKey).(io.Writer)
}

func ErrWithContext(ctx context.Context, err io.Writer) context.Context {
	return context.WithValue(ctx, errKey, err)
}
