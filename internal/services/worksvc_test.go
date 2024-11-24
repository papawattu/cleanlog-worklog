package services_test

import (
	"context"
	"testing"
	"time"

	common "github.com/papawattu/cleanlog-common"
	repo "github.com/papawattu/cleanlog-common"
	"github.com/papawattu/cleanlog-worklog/internal/models"
	"github.com/papawattu/cleanlog-worklog/internal/services"
)

func TestWorkServiceImp_CreateWorkLog(t *testing.T) {

	type fields struct {
		ctx  context.Context
		repo repo.Repository[*models.WorkLog, string]
	}
	type args struct {
		ctx         context.Context
		description string
		date        time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Create work log",
			fields: fields{
				ctx:  context.Background(),
				repo: common.NewInMemoryRepository[*models.WorkLog](),
			},
			args: args{
				ctx:         context.Background(),
				description: "Test work log",
				date:        time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			wsi := services.NewWorkService(tt.fields.ctx, tt.fields.repo)

			got, err := wsi.CreateWorkLog(tt.args.ctx, tt.args.description, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("WorkServiceImp.CreateWorkLog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("WorkServiceImp.CreateWorkLog() = %v", got)
			wl, err := wsi.GetWorkLog(tt.args.ctx, got)
			if err != nil {
				t.Errorf("WorkServiceImp.GetWorkLog() error = %v", err)
			}
			if wl == nil {
				t.Errorf("WorkServiceImp.GetWorkLog() = %v, want %v", wl, tt.want)
			}
			if wl.WorkLogDescription != tt.args.description {
				t.Errorf("WorkServiceImp.GetWorkLog().WorkLogDescription = %v, want %v", wl.WorkLogDescription, tt.args.description)
			}
			if wl.WorkLogDate != tt.args.date {
				t.Errorf("WorkServiceImp.GetWorkLog().WorkLogDate = %v, want %v", wl.WorkLogDate, tt.args.date)
			}

		})
	}

}
