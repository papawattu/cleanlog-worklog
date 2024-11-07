package services

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/papawattu/cleanlog-worklog/internal/models"
	"github.com/papawattu/cleanlog-worklog/internal/repo"
)

type WorkService interface {
	CreateWorkLog(ctx context.Context, description string, date time.Time) (int, error)

	DeleteWorkLog(ctx context.Context, id int) error

	GetWorkLog(ctx context.Context, id int) (*models.WorkLog, error)

	GetAllWorkLog(ctx context.Context, user int) ([]*models.WorkLog, error)

	UpdateWorkLog(ctx context.Context, id int, description string, date time.Time) error
}

type WorkServiceImp struct {
	ctx  context.Context
	repo repo.Repository[*models.WorkLog, int]
}

func nextId() int {
	return rand.Intn(1000)
}

func (wsi *WorkServiceImp) CreateWorkLog(ctx context.Context, description string, date time.Time) (int, error) {

	wl, err := models.NewWorkLog(description, date)
	if err != nil {
		log.Fatalf("Error starting work: %v", err)
		return 0, err
	}

	nextId := nextId()
	wl.WorkLogID = &nextId

	err = wsi.repo.Save(ctx, &wl)
	if err != nil {
		log.Fatalf("Error saving work log: %v", err)
		return 0, err
	}
	return nextId, nil
}

func (wsi *WorkServiceImp) LogWork(ctx context.Context, id int, t models.Task) error {

	wl, err := wsi.repo.Get(ctx, id)
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

func (wsi *WorkServiceImp) DeleteWorkLog(ctx context.Context, id int) error {

	wl, err := wsi.repo.Get(ctx, id)
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

	err = wsi.repo.Delete(ctx, wl)
	if err != nil {
		log.Fatalf("Error deleting work log: %v", err)
		return err
	}
	return nil
}

func (wsi *WorkServiceImp) GetWorkLog(ctx context.Context, id int) (*models.WorkLog, error) {

	wl, err := wsi.repo.Get(ctx, id)

	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
		return nil, err
	}

	return wl, nil
}
func (wsi *WorkServiceImp) GetAllWorkLog(ctx context.Context, user int) ([]*models.WorkLog, error) {

	wls, err := wsi.repo.GetAll(ctx)

	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
		return nil, err
	}

	return wls, nil
}

func (wsi *WorkServiceImp) UpdateWorkLog(ctx context.Context, id int, description string, date time.Time) error {

	wl, err := wsi.repo.Get(ctx, id)
	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
		return err
	}

	wl.WorkLogDescription = description
	wl.WorkLogDate = date

	err = wsi.repo.Save(ctx, wl)
	if err != nil {
		log.Fatalf("Error updating work log: %v", err)
		return err
	}

	return nil
}
func NewWorkService(ctx context.Context, repo repo.Repository[*models.WorkLog, int]) WorkService {

	return &WorkServiceImp{
		ctx:  ctx,
		repo: repo,
	}
}
