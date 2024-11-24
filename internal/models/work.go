package models

import (
	"strconv"
	"time"

	common "github.com/papawattu/cleanlog-common"
)

type Task struct {
	TaskID int
}
type WorkLog struct {
	common.BaseEntity[int]
	WorkLogID          *int
	WorkLogDate        time.Time
	WorkLogTimeInSecs  int
	WorkLogDescription string
	Tasks              []Task
	UserID             int
}
type Work interface {
	LogWork(WorkLog) error
}

func NewWorkLog(description string, date time.Time) (WorkLog, error) {

	wl := WorkLog{
		WorkLogID:          nil,
		WorkLogDate:        date,
		WorkLogDescription: description,
		Tasks:              make([]Task, 0),
		UserID:             0,
	}
	return wl, nil
}

func (wl *WorkLog) LogWork(t Task) error {
	wl.Tasks = append(wl.Tasks, t)
	return nil
}

func (wl *WorkLog) EndWork() error {
	now := time.Now()
	worked := now.Sub(wl.WorkLogDate)
	wl.WorkLogTimeInSecs = int(worked.Seconds())
	wl.LastUpdateDate = now
	return nil
}

func (wl *WorkLog) AddTask(t Task) error {
	wl.Tasks = append(wl.Tasks, t)
	return nil
}

func (wl *WorkLog) RemoveTask(t Task) error {
	for i, task := range wl.Tasks {
		if task.TaskID == t.TaskID {
			wl.Tasks = append(wl.Tasks[:i], wl.Tasks[i+1:]...)
		}
	}
	return nil
}

func (wl *WorkLog) HasTask(t Task) bool {
	for _, task := range wl.Tasks {
		if task.TaskID == t.TaskID {
			return true
		}
	}
	return false
}

func (wl *WorkLog) ChangeDescription(description string) error {
	wl.WorkLogDescription = description
	return nil
}

func (wl *WorkLog) ChangeDate(date time.Time) error {
	wl.WorkLogDate = date
	return nil
}

func (wl *WorkLog) ChangeUserID(id int) error {

	wl.UserID = id
	return nil
}

func (wl *WorkLog) GetID() string {
	return strconv.Itoa(*wl.WorkLogID)
}
