package repo

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"strconv"
	"time"

	repo "github.com/papawattu/cleanlog-eventstore/repository"
	"github.com/papawattu/cleanlog-worklog/internal/models"
)

type InMemoryWorkLogRepository struct {
	WorkLogs map[int]*models.WorkLog
}

func (wri *InMemoryWorkLogRepository) Save(ctx context.Context, wl *models.WorkLog) error {
	if wl.WorkLogID == nil {
		return errors.New("work log ID is required")
	}

	// if _, ok := wri.WorkLogs[*wl.WorkLogID]; ok {
	// 	return errors.New("work log already exists")
	// }

	wl.LastUpdateDate = time.Now()
	wri.WorkLogs[*wl.WorkLogID] = wl

	log.Printf("Work log saved: %v", wl)
	return nil
}

func (wri *InMemoryWorkLogRepository) Get(ctx context.Context, id int) (*models.WorkLog, error) {
	wl, ok := wri.WorkLogs[id]
	if !ok {
		return nil, nil
	}
	return wl, nil
}

func (wri *InMemoryWorkLogRepository) GetAll(ctx context.Context) ([]*models.WorkLog, error) {
	userID, ok := ctx.Value("user").(int)
	if !ok {
		return nil, errors.New("user ID not found in context")
	}

	workLogs := []*models.WorkLog{}
	for _, wl := range wri.WorkLogs {
		if wl.UserID == userID {
			workLogs = append(workLogs, wl)
		}
	}
	return workLogs, nil
}

func (wri *InMemoryWorkLogRepository) Delete(ctx context.Context, wl *models.WorkLog) error {

	if wl.WorkLogID == nil {
		return errors.New("work log ID is required")
	}

	if _, ok := wri.WorkLogs[*wl.WorkLogID]; !ok {
		slog.Error("Work log with ID %s not found", strconv.Itoa(*wl.WorkLogID), nil)
		return errors.New("work log not found " + strconv.Itoa(*wl.WorkLogID))
	}
	delete(wri.WorkLogs, *wl.WorkLogID)
	return nil
}

func (wri *InMemoryWorkLogRepository) Update(ctx context.Context, wl *models.WorkLog) error {
	if wl.WorkLogID == nil {
		return errors.New("work log ID is required")
	}

	if _, ok := wri.WorkLogs[*wl.WorkLogID]; !ok {
		return errors.New("work log not found")
	}

	wl.LastUpdateDate = time.Now()
	wri.WorkLogs[*wl.WorkLogID] = wl

	return nil
}

func (wri *InMemoryWorkLogRepository) GetId(ctx context.Context, wl *models.WorkLog) (int, error) {

	if _, ok := wri.WorkLogs[*wl.WorkLogID]; !ok {
		return 0, errors.New("work log not found")
	}
	return *wl.WorkLogID, nil
}

func (wri *InMemoryWorkLogRepository) Exists(ctx context.Context, id int) (bool, error) {

	_, ok := wri.WorkLogs[id]
	return ok, nil
}
func NewWorkLogRepository() repo.Repository[*models.WorkLog, int] {
	return &InMemoryWorkLogRepository{
		WorkLogs: make(map[int]*models.WorkLog),
	}
}
