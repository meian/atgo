package io

import "io"

type Writer = io.Writer
type Reader = io.Reader

var Discard = io.Discard

func Copy(dst Writer, src Reader) (written int64, err error) {
	return io.Copy(dst, src)
}
