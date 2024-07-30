package requests_test

import (
	"testing"

	"github.com/meian/atgo/crawler/requests"
	"github.com/stretchr/testify/assert"
)

func TestTask_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     requests.Task
		wantErr bool
	}{
		{
			name:    "success",
			req:     requests.Task{ContestID: "abc123", TaskID: "abc123_a"},
			wantErr: false,
		},
		{
			name:    "no contest id",
			req:     requests.Task{ContestID: "", TaskID: "abc123_a"},
			wantErr: true,
		},
		{
			name:    "no task id",
			req:     requests.Task{ContestID: "abc123", TaskID: ""},
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
