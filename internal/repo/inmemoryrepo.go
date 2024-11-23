package repo

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	common "github.com/papawattu/cleanlog-common"
	"github.com/papawattu/cleanlog-worklog/internal/models"
)

type InMemoryWorkLogRepository struct {
	WorkLogs map[int]*models.WorkLog
}

func (wri *InMemoryWorkLogRepository) Create(ctx context.Context, wl *models.WorkLog) error {
	if wl.WorkLogID == nil {
		return errors.New("work log ID is required")
	}

	if _, ok := wri.WorkLogs[*wl.WorkLogID]; ok {
		return errors.New("work log already exists")
	}

	wl.LastUpdateDate = time.Now()
	wri.WorkLogs[*wl.WorkLogID] = wl

	slog.Info("Work log created in repository", "Worklog", wl)
	return nil
}
func (wri *InMemoryWorkLogRepository) Save(ctx context.Context, wl *models.WorkLog) error {
	if wl.WorkLogID == nil {
		return errors.New("work log ID is required")
	}

	wl.LastUpdateDate = time.Now()
	wri.WorkLogs[*wl.WorkLogID] = wl

	slog.Info("Work log saved to repository", "Worklog", wl)
	return nil
}

func (wri *InMemoryWorkLogRepository) Get(ctx context.Context, id string) (*models.WorkLog, error) {
	i, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	wl, ok := wri.WorkLogs[i]
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

func (wri *InMemoryWorkLogRepository) GetId(ctx context.Context, wl *models.WorkLog) (string, error) {

	if _, ok := wri.WorkLogs[*wl.WorkLogID]; !ok {
		return "", errors.New("work log not found")
	}
	i := strconv.Itoa(*wl.WorkLogID)
	return i, nil
}

func (wri *InMemoryWorkLogRepository) Exists(ctx context.Context, id string) (bool, error) {

	i, err := strconv.Atoi(id)
	if err != nil {
		return false, err
	}
	_, ok := wri.WorkLogs[i]
	return ok, nil
}
func NewWorkLogRepository() common.Repository[*models.WorkLog, string] {
	return &InMemoryWorkLogRepository{
		WorkLogs: make(map[int]*models.WorkLog),
	}
}
