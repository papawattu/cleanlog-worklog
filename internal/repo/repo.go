package repo

import (
	"time"

	"github.com/papawattu/cleanlog-worklog/internal/models"
)

type WorkLogRepository interface {
	SaveWorkLog(wl *models.WorkLog) error
	GetWorkLog(id int) (*models.WorkLog, error)
	GetAllWorkLogsForUser(userID int) ([]*models.WorkLog, error)
	DeleteWorkLog(id int) error
}

type WorkLogRepositoryImp struct {
	WorkLogs map[int]*models.WorkLog
}

func (wri *WorkLogRepositoryImp) SaveWorkLog(wl *models.WorkLog) error {
	wl.LastUpdateDate = time.Now()
	wri.WorkLogs[*wl.WorkLogID] = wl
	return nil
}

func (wri *WorkLogRepositoryImp) GetWorkLog(id int) (*models.WorkLog, error) {
	wl, ok := wri.WorkLogs[id]
	if !ok {
		return nil, nil
	}
	return wl, nil
}

func (wri *WorkLogRepositoryImp) GetAllWorkLogsForUser(userID int) ([]*models.WorkLog, error) {
	workLogs := []*models.WorkLog{}
	for _, wl := range wri.WorkLogs {
		if wl.UserID == userID {
			workLogs = append(workLogs, wl)
		}
	}
	return workLogs, nil
}

func (wri *WorkLogRepositoryImp) DeleteWorkLog(id int) error {
	delete(wri.WorkLogs, id)
	return nil
}
func NewWorkLogRepository() WorkLogRepository {
	return &WorkLogRepositoryImp{
		WorkLogs: make(map[int]*models.WorkLog),
	}
}
