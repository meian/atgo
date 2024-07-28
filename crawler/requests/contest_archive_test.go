package requests_test

import (
	"testing"

	"github.com/meian/atgo/constant"
	"github.com/meian/atgo/crawler/requests"
	"github.com/meian/atgo/util"
	"github.com/stretchr/testify/assert"
)

func TestContestArchive_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     requests.ContestArchive
		wantErr bool
	}{
		{
			name: "full params",
			req: requests.ContestArchive{
				Page:      1,
				RatedType: util.ToPtr(constant.RatedTypeABC),
				Category:  util.ToPtr(constant.CategoryDailyTraining),
				Keyword:   util.ToPtr("keyword"),
			},
			wantErr: false,
		},
		{
			name: "0 page",
			req: requests.ContestArchive{
				Page:      0,
				RatedType: util.ToPtr(constant.RatedTypeABC),
				Category:  util.ToPtr(constant.CategoryDailyTraining),
				Keyword:   util.ToPtr("keyword"),
			},
			wantErr: true,
		},
		{
			name: "minus page",
			req: requests.ContestArchive{
				Page:      -1,
				RatedType: util.ToPtr(constant.RatedTypeABC),
				Category:  util.ToPtr(constant.CategoryDailyTraining),
				Keyword:   util.ToPtr("keyword"),
			},
			wantErr: true,
		},
		{
			name: "nil rated type",
			req: requests.ContestArchive{
				Page:      1,
				RatedType: nil,
				Category:  util.ToPtr(constant.CategoryDailyTraining),
				Keyword:   util.ToPtr("keyword"),
			},
			wantErr: false,
		},
		{
			name: "invalid rated type",
			req: requests.ContestArchive{
				Page:      1,
				RatedType: util.ToPtr(constant.RatedType(-1)),
				Category:  util.ToPtr(constant.CategoryDailyTraining),
				Keyword:   util.ToPtr("keyword"),
			},
			wantErr: true,
		},
		{
			name: "nil category",
			req: requests.ContestArchive{
				Page:      1,
				RatedType: util.ToPtr(constant.RatedTypeABC),
				Category:  nil,
				Keyword:   util.ToPtr("keyword"),
			},
			wantErr: false,
		},
		{
			name: "invalid category",
			req: requests.ContestArchive{
				Page:      1,
				RatedType: util.ToPtr(constant.RatedTypeABC),
				Category:  util.ToPtr(constant.ContestCategory(-1)),
				Keyword:   util.ToPtr("keyword"),
			},
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
