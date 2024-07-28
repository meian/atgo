package requests_test

import (
	"testing"

	"github.com/meian/atgo/crawler/requests"
	"github.com/stretchr/testify/assert"
)

func TestContestTask_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     requests.ContestTask
		wantErr bool
	}{
		{
			name:    "success",
			req:     requests.ContestTask{ContestID: "abc123"},
			wantErr: false,
		},
		{
			name:    "no contest id",
			req:     requests.ContestTask{ContestID: ""},
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
