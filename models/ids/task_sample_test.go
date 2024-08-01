package ids_test

import (
	"testing"

	"github.com/meian/atgo/models/ids"
	"github.com/stretchr/testify/assert"
)

func TestTaskSampleID_Validate(t *testing.T) {
	tests := []struct {
		name     string
		id       ids.TaskSampleID
		errMsgPt string
	}{
		{
			name:     "success",
			id:       "abc123_a__1",
			errMsgPt: "",
		},
		{
			name:     "not contains _",
			id:       "abc123a__1",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "not contains __",
			id:       "abc123_a1",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "_ instead of __",
			id:       "abc123_a_1",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "empty",
			id:       "",
			errMsgPt: `empty id`,
		},
		{
			name:     "64 chars",
			id:       "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvw_a__1",
			errMsgPt: "",
		},
		{
			name:     "65 chars",
			id:       "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwx_a__1",
			errMsgPt: `is too long`,
		},
		{
			name:     "first char is -",
			id:       "-abc123_a__1",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "last char is -",
			id:       "abc123_a__1-",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "contains -",
			id:       "abc-123_a__1",
			errMsgPt: "",
		},
		{
			name:     "double -",
			id:       "abc--123_a__1",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "first char is _",
			id:       "_abc123_a__1",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "last char is _",
			id:       "abc123_a__1_",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "contains _",
			id:       "abc_123_a__1",
			errMsgPt: "",
		},
		{
			name:     "double _",
			id:       "abc__123_a__1",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "contains double - and _",
			id:       "abc-_123_a__1",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "contains symbol",
			id:       "abc!123_a__1",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "contains space",
			id:       "abc 123_a__1",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "contains tab",
			id:       "abc\t123_a__1",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "contains newline",
			id:       "abc\n123_a__1",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "contains carriage return",
			id:       "abc\r123_a__1",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "multi byte",
			id:       "ａｂｃ１２３_ａ__１",
			errMsgPt: `invalid .+ format:`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			err := tt.id.Validate()
			if tt.errMsgPt != "" {
				if !assert.Error(err) {
					return
				}
				assert.Regexp(tt.errMsgPt, err.Error())
				return
			}
			assert.NoError(err)
		})
	}
}
