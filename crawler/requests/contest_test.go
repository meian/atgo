package requests_test

import (
	"testing"

	"github.com/meian/atgo/crawler/requests"
	"github.com/stretchr/testify/assert"
)

func TestContest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     requests.Contest
		wantErr bool
	}{
		{
			name:    "success",
			req:     requests.Contest{ContestID: "abc123"},
			wantErr: false,
		},
		{
			name:    "no contest id",
			req:     requests.Contest{ContestID: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
