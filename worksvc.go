package main

import (
	"log"
	"math/rand"
)

type WorkService interface {
	StartWork(description string) (int, error)

	LogWork(id int, t Task) error

	EndWork(id int) error

	GetWorkLog(id int) (*WorkLog, error)
}

type WorkServiceImp struct {
	repo WorkLogRepository
}

func nextId() int {
	return rand.Intn(1000)
}

func (wsi *WorkServiceImp) StartWork(description string) (int, error) {

	wl, err := NewWorkLog(description)
	if err != nil {
		log.Fatalf("Error starting work: %v", err)
		return 0, err
	}

	nextId := nextId()
	wl.WorkLogID = &nextId

	wsi.repo.SaveWorkLog(&wl)

	err = wsi.repo.SaveWorkLog(&wl)
	if err != nil {
		log.Fatalf("Error saving work log: %v", err)
		return 0, err
	}
	return nextId, nil
}

func (wsi *WorkServiceImp) LogWork(id int, t Task) error {

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

func (wsi *WorkServiceImp) EndWork(id int) error {

	wl, err := wsi.repo.GetWorkLog(id)
	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
		return err
	}

	err = wl.EndWork()
	if err != nil {
		log.Fatalf("Error ending work")
		return err
	}

	return nil
}

func (wsi *WorkServiceImp) GetWorkLog(id int) (*WorkLog, error) {

	wl, err := wsi.repo.GetWorkLog(id)

	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
		return nil, err
	}

	return wl, nil
}
func NewWorkService(repo WorkLogRepository) WorkService {

	return &WorkServiceImp{
		repo: repo,
	}
}
