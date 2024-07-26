package io

import (
	"bytes"
	"io"
)

func WithReadAction(r io.Reader, f func(io.Reader) error) (io.Reader, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	if err := f(tee); err != nil {
		return nil, err
	}
	if _, err := io.Copy(io.Discard, tee); err != nil {
		return nil, err
	}
	return &buf, nil
}
