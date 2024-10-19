package main

type WorkLogRepository interface {
	SaveWorkLog(wl *WorkLog) error
	GetWorkLog(id int) (*WorkLog, error)
}

type WorkLogRepositoryImp struct {
	WorkLogs map[int]*WorkLog
}

func (wri *WorkLogRepositoryImp) SaveWorkLog(wl *WorkLog) error {
	wri.WorkLogs[*wl.WorkLogID] = wl
	return nil
}

func (wri *WorkLogRepositoryImp) GetWorkLog(id int) (*WorkLog, error) {
	wl, ok := wri.WorkLogs[id]
	if !ok {
		return nil, nil
	}
	return wl, nil
}
func NewWorkLogRepository() WorkLogRepository {
	return &WorkLogRepositoryImp{
		WorkLogs: make(map[int]*WorkLog),
	}
}
