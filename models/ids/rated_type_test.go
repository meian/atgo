package ids_test

import (
	"testing"

	"github.com/meian/atgo/models/ids"
	"github.com/stretchr/testify/assert"
)

func TestRatedType_Validate(t *testing.T) {
	tests := []struct {
		name     string
		rt       ids.RatedType
		errMsgPt string
	}{
		{
			name:     "abc",
			rt:       "abc",
			errMsgPt: "",
		},
		{
			name:     "arc",
			rt:       "arc",
			errMsgPt: "",
		},
		{
			name:     "agc",
			rt:       "agc",
			errMsgPt: "",
		},
		{
			name:     "ahc",
			rt:       "ahc",
			errMsgPt: "",
		},
		{
			name:     "not defined",
			rt:       "xyz",
			errMsgPt: `invalid .+ format: xyz`,
		},
		{
			name:     "upper case",
			rt:       "ABC",
			errMsgPt: `invalid .+ format: ABC`,
		},
		{
			name:     "upper first",
			rt:       "Abc",
			errMsgPt: `invalid .+ format: Abc`,
		},
		{
			name:     "empty",
			rt:       "",
			errMsgPt: `empty id`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			err := tt.rt.Validate()
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
