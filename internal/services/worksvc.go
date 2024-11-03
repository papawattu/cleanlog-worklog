package services

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/papawattu/cleanlog-worklog/internal/models"
	"github.com/papawattu/cleanlog-worklog/internal/repo"
)

type WorkService interface {
	CreateWorkLog(description string, date time.Time) (int, error)

	//	LogWork(id int, t Task) error

	DeleteWorkLog(id int) error

	GetWorkLog(id int) (*models.WorkLog, error)

	GetAllWorkLog(user int) ([]*models.WorkLog, error)

	UpdateWorkLog(id int, description string, date time.Time) error
}

type WorkServiceImp struct {
	repo repo.WorkLogRepository
}

func nextId() int {
	return rand.Intn(1000)
}

func (wsi *WorkServiceImp) CreateWorkLog(description string, date time.Time) (int, error) {

	wl, err := models.NewWorkLog(description, date)
	if err != nil {
		log.Fatalf("Error starting work: %v", err)
		return 0, err
	}

	nextId := nextId()
	wl.WorkLogID = &nextId

	err = wsi.repo.SaveWorkLog(&wl)
	if err != nil {
		log.Fatalf("Error saving work log: %v", err)
		return 0, err
	}
	return nextId, nil
}

func (wsi *WorkServiceImp) LogWork(id int, t models.Task) error {

	wl, err := wsi.repo.GetWorkLog(id)
	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
		return err
	}

	err = wl.LogWork(t)
	if err != nil {
		log.Fatalf("Error logging work: %v", err)
		return err
	}

	return nil
}

func (wsi *WorkServiceImp) DeleteWorkLog(id int) error {

	wl, err := wsi.repo.GetWorkLog(id)
	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
		return err
	}

	if wl == nil {
		return errors.New("Work log not found")
	}

	err = wl.EndWork()
	if err != nil {
		log.Fatalf("Error ending work")
		return err
	}

	err = wsi.repo.DeleteWorkLog(id)
	if err != nil {
		log.Fatalf("Error deleting work log: %v", err)
		return err
	}
	return nil
}

func (wsi *WorkServiceImp) GetWorkLog(id int) (*models.WorkLog, error) {

	wl, err := wsi.repo.GetWorkLog(id)

	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
		return nil, err
	}

	return wl, nil
}
func (wsi *WorkServiceImp) GetAllWorkLog(user int) ([]*models.WorkLog, error) {

	wls, err := wsi.repo.GetAllWorkLogsForUser(user)

	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
		return nil, err
	}

	return wls, nil
}

func (wsi *WorkServiceImp) UpdateWorkLog(id int, description string, date time.Time) error {

	wl, err := wsi.repo.GetWorkLog(id)
	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
		return err
	}

	wl.WorkLogDescription = description
	wl.WorkLogDate = date

	err = wsi.repo.SaveWorkLog(wl)
	if err != nil {
		log.Fatalf("Error updating work log: %v", err)
		return err
	}

	return nil
}
func NewWorkService(repo repo.WorkLogRepository) WorkService {

	return &WorkServiceImp{
		repo: repo,
	}
}
