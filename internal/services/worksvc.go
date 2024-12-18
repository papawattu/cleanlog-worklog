package services

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"math/rand"
	"strconv"
	"time"

	repo "github.com/papawattu/cleanlog-common"
	"github.com/papawattu/cleanlog-worklog/internal/models"
)

type WorkService interface {
	CreateWorkLog(ctx context.Context, description string, date time.Time) (int, error)

	DeleteWorkLog(ctx context.Context, id int) error

	GetWorkLog(ctx context.Context, id int) (*models.WorkLog, error)

	GetAllWorkLog(ctx context.Context, user int) ([]*models.WorkLog, error)

	UpdateWorkLog(ctx context.Context, id int, description string, date time.Time) error

	AddTaskToWorkLog(ctx context.Context, id int, t models.Task) error

	RemoveTaskFromWorkLog(ctx context.Context, id int, t models.Task) error
}

type WorkServiceImp struct {
	ctx  context.Context
	repo repo.Repository[*models.WorkLog, string]
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
	slog.Info("Creating work log", "id", nextId)
	err = wsi.repo.Create(ctx, &wl)
	if err != nil {
		slog.Error("Error saving work log", "error", err)
		return nextId, err
	}
	return nextId, nil
}

func (wsi *WorkServiceImp) LogWork(ctx context.Context, id int, t models.Task) error {

	wl, err := wsi.repo.Get(ctx, strconv.Itoa(id))
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

	wl, err := wsi.repo.Get(ctx, strconv.Itoa(id))
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

	wl, err := wsi.repo.Get(ctx, strconv.Itoa(id))

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

	wl, err := wsi.repo.Get(ctx, strconv.Itoa(id))
	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
		return err
	}

	if wl == nil {
		return errors.New("Work log not found")
	}

	if description != "" {
		wl.WorkLogDescription = description
	}

	if date != (time.Time{}) {
		wl.WorkLogDate = date
	}

	err = wsi.repo.Save(ctx, wl)
	if err != nil {
		log.Fatalf("Error updating work log: %v", err)
		return err
	}

	return nil
}

func (wsi *WorkServiceImp) AddTaskToWorkLog(ctx context.Context, id int, t models.Task) error {

	wl, err := wsi.repo.Get(ctx, strconv.Itoa(id))
	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
		return err
	}

	err = wl.AddTask(t)
	if err != nil {
		log.Fatalf("Error adding task: %v", err)
		return err
	}

	err = wsi.repo.Save(ctx, wl)
	if err != nil {
		log.Fatalf("Error saving work log: %v", err)
		return err
	}

	return nil
}

func (wsi *WorkServiceImp) RemoveTaskFromWorkLog(ctx context.Context, id int, t models.Task) error {

	wl, err := wsi.repo.Get(ctx, strconv.Itoa(id))
	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
		return err
	}

	if wl == nil {
		return errors.New("Work log not found")
	}

	err = wl.RemoveTask(t)
	if err != nil {
		log.Fatalf("Error removing task: %v", err)
		return err
	}

	err = wsi.repo.Save(ctx, wl)
	if err != nil {
		log.Fatalf("Error saving work log: %v", err)
		return err
	}

	return nil
}
func NewWorkService(ctx context.Context, repo repo.Repository[*models.WorkLog, string]) WorkService {

	return &WorkServiceImp{
		ctx:  ctx,
		repo: repo,
	}
}
