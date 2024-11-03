package models

import "time"

type Task struct {
	TaskID int
}
type WorkLog struct {
	WorkLogID          *int
	WorkLogDate        time.Time
	WorkLogTimeInSecs  int
	WorkLogDescription string
	Tasks              []Task
	UserID             int
	CreationDate       time.Time
	LastUpdateDate     time.Time
}
type Work interface {
	LogWork(WorkLog) error
}

func NewWorkLog(description string, date time.Time) (WorkLog, error) {

	wl := WorkLog{
		WorkLogID:          nil,
		WorkLogDate:        date,
		WorkLogDescription: description,
		CreationDate:       time.Now(),
		LastUpdateDate:     time.Now(),
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

func (wl *WorkLog) ChangeDescription(description string) error {
	wl.WorkLogDescription = description
	return nil
}
