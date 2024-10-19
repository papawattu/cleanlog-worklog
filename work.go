package main

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
}
type Work interface {
	LogWork(WorkLog) error
}

func NewWorkLog(description string) (WorkLog, error) {

	wl := WorkLog{
		WorkLogID:          nil,
		WorkLogDate:        time.Now(),
		WorkLogDescription: description,
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
	return nil
}
