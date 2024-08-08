package ids_test

import (
	"testing"

	"github.com/meian/atgo/models/ids"
	"github.com/stretchr/testify/assert"
)

func TestContestTaskID_Validate(t *testing.T) {
	tests := []struct {
		name     string
		id       ids.ContestTaskID
		errMsgPt string
	}{
		{
			name:     "success",
			id:       "abc123--abc123_a",
			errMsgPt: "",
		},
		{
			name:     "not contains --",
			id:       "abc123-abc123_a",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "not contains _",
			id:       "abc123--abc123a",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "empty",
			id:       "",
			errMsgPt: `empty id`,
		},
		{
			name:     "64 chars",
			id:       "0123456789--ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxy",
			errMsgPt: "",
		},
		{
			name:     "65 chars",
			id:       "0123456789--ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz",
			errMsgPt: `is too long`,
		},
		{
			name:     "first char is -",
			id:       "-abc123--abc123_a",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "last char is -",
			id:       "abc123--abc123_a-",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "contains -",
			id:       "abc-123--abc123_a",
			errMsgPt: "",
		},
		{
			name:     "double -",
			id:       "abc--123--abc123_a",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "first char is _",
			id:       "_abc123--abc123_a",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "last char is _",
			id:       "abc123--abc123_a_",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "contains _",
			id:       "abc_123--abc123_a",
			errMsgPt: "",
		},
		{
			name:     "double _",
			id:       "abc__123--abc123_a",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "contains double - and _",
			id:       "abc--123--abc123_a",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "contains symbol",
			id:       "abc!123--abc123_a",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "contains space",
			id:       "abc 123--abc123_a",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "contains tab",
			id:       "abc\t123--abc123_a",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "contains newline",
			id:       "abc\n123--abc123_a",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "contains carriage return",
			id:       "abc\r123--abc123_a",
			errMsgPt: `invalid .+ format:`,
		},
		{
			name:     "multi byte",
			id:       "ａｂｃ１２３--ａｂｃ１２３_ａ",
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