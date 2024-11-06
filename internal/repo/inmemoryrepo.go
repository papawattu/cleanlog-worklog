package repo

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/papawattu/cleanlog-worklog/internal/models"
)

type InMemoryWorkLogRepository struct {
	WorkLogs map[int]*models.WorkLog
}

func (wri *InMemoryWorkLogRepository) Save(ctx context.Context, wl *models.WorkLog) error {
	if wl.WorkLogID == nil {
		return errors.New("work log ID is required")
	}

	if _, ok := wri.WorkLogs[*wl.WorkLogID]; ok {
		return errors.New("work log already exists")
	}

	wl.LastUpdateDate = time.Now()
	wri.WorkLogs[*wl.WorkLogID] = wl
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
	userID := ctx.Value("userID")

	workLogs := []*models.WorkLog{}
	for _, wl := range wri.WorkLogs {
		if wl.UserID == userID {
			workLogs = append(workLogs, wl)
		}
	}
	return workLogs, nil
}

func (wri *InMemoryWorkLogRepository) Delete(ctx context.Context, id int) error {
	if _, ok := wri.WorkLogs[id]; !ok {
		log.Printf("Work log with ID %d not found", id)
		return errors.New("work log not found " + strconv.Itoa(id))
	}
	delete(wri.WorkLogs, id)
	return nil
}
func NewWorkLogRepository() Repository[*models.WorkLog, int] {
	return &InMemoryWorkLogRepository{
		WorkLogs: make(map[int]*models.WorkLog),
	}
}
