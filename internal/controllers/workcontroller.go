package controllers

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/papawattu/cleanlog-worklog/internal/models"
	"github.com/papawattu/cleanlog-worklog/internal/services"
	"github.com/papawattu/cleanlog-worklog/types"
)

type ControllerPaths map[string]func(wc *WorkController, ctx context.Context) func(http.ResponseWriter, *http.Request)

func (cp ControllerPaths) GetPaths() map[string]func(wc *WorkController, ctx context.Context) func(http.ResponseWriter, *http.Request) {
	return cp
}

type WorkController struct {
	workService services.WorkService
	server      *http.ServeMux
	controllers ControllerPaths
}

func inlineTasks(tasks []models.Task) []int {
	if tasks == nil {
		return make([]int, 0)
	}
	t := make([]int, 0)
	for _, task := range tasks {
		t = append(t, task.TaskID)
	}
	return t
}
func (wc *WorkController) PostRequest(ctx context.Context) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Creating work log")
		var t types.CreateWorkRequest

		json.NewDecoder(r.Body).Decode(&t)

		startDate, err := time.Parse("2006-01-02", t.Date)

		if err != nil {
			slog.Error("Invalid date format %v", err)
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}

		workID, err := wc.workService.CreateWorkLog(ctx, t.Description, startDate)
		if err != nil {
			slog.Error("Error starting work: %v", err)
		}

		w.Header().Set("Location", "/api/worklog/"+strconv.Itoa(workID))
		w.WriteHeader(http.StatusCreated)
	}

}

func (wc *WorkController) PatchRequest(ctx context.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		workId := r.PathValue("workid")
		if workId == "" {
			http.Error(w, "workId is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(workId)
		if err != nil {
			http.Error(w, "workId must be an integer", http.StatusBadRequest)
			return
		}

		slog.Debug("Updating work log by id", slog.Int("id", id))

		var t types.UpdateWorkRequest

		json.NewDecoder(r.Body).Decode(&t)

		var startDate time.Time
		if t.Date != "" {
			startDate, err = time.Parse("2006-01-02", t.Date)

			if err != nil {
				slog.Error("Invalid date format", "Date", t.Date)
				http.Error(w, "Invalid date format", http.StatusBadRequest)
				return
			}
		}

		err = wc.workService.UpdateWorkLog(ctx, id, t.Description, startDate, t.TaskIds)
		if err != nil {
			http.Error(w, "Error updating work", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (wc *WorkController) GetRequestById(ctx context.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		workId := r.PathValue("workid")
		if workId == "" {
			http.Error(w, "workId is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(workId)
		if err != nil {
			http.Error(w, "workId must be an integer", http.StatusBadRequest)
			return
		}

		log.Printf("Getting work log by id %d", id)

		work, err := wc.workService.GetWorkLog(ctx, id)
		if err != nil {
			http.Error(w, "Error getting work", http.StatusNotFound)
			return
		}

		if work == nil {
			http.Error(w, "Work log not found", http.StatusNotFound)
			return
		}

		log.Printf("Work log: %v", work)
		wlr := types.WorkResponse{
			WorkID:      *work.WorkLogID,
			Description: work.WorkLogDescription,
			TaskIds:     inlineTasks(work.Tasks),
			Date:        work.WorkLogDate.Format("2006-01-02"),
			CreatedAt:   work.CreationDate.Format(time.RFC3339Nano),
			UpdatedAt:   work.LastUpdateDate.Format(time.RFC3339Nano),
			UserID:      work.UserID,
		}
		json.NewEncoder(w).Encode(wlr)
	}
}

func (wc *WorkController) GetRequestAll(ctx context.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Debug("Getting all work logs")

		// TODO: Implement user id

		user, ok := r.Context().Value("user").(int)

		if !ok {
			slog.Error("User ID not found in context")
			http.Error(w, "User ID not found in context", http.StatusNotFound)
			return
		}
		workLogs, err := wc.workService.GetAllWorkLog(r.Context(), user)
		if err != nil {
			slog.Error("Error getting work logs: %v", err)
			http.Error(w, "Error getting work logs", http.StatusNotFound)
			return
		}

		wlr := &types.ListWorkResponse{}

		wlr.WorkResponses = make([]types.WorkResponse, 0)

		for _, workLog := range workLogs {
			wlr.WorkResponses = append(wlr.WorkResponses,
				types.WorkResponse{
					WorkID:      *workLog.WorkLogID,
					Description: workLog.WorkLogDescription,
					TaskIds:     inlineTasks(workLog.Tasks),
					Date:        workLog.WorkLogDate.Format("2006-01-02"),
					CreatedAt:   workLog.CreationDate.Format("2006-01-02"),
					UpdatedAt:   workLog.LastUpdateDate.Format("2006-01-02"),
				})
		}
		json.NewEncoder(w).Encode(wlr)
	}
}

func (wc *WorkController) DeleteRequest(ctx context.Context) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		workId := r.PathValue("workid")
		if workId == "" {
			slog.Error("workId is required")
			http.Error(w, "workId is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(workId)
		if err != nil {
			slog.Error("workId must be an integer")
			http.Error(w, "workId must be an integer", http.StatusBadRequest)
			return
		}

		slog.Debug("Deleting work log by id", slog.Int("id", id))

		err = wc.workService.DeleteWorkLog(ctx, id)
		if err != nil {
			http.Error(w, "Error deleting work", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
func NewWorkController(ctx context.Context, server *http.ServeMux,
	workService services.WorkService) *WorkController {

	wc := &WorkController{
		workService: workService,
		controllers: ControllerPaths{
			"POST /worklog":            (*WorkController).PostRequest,
			"GET /worklog/{workid}":    (*WorkController).GetRequestById,
			"GET /worklog/":            (*WorkController).GetRequestAll,
			"DELETE /worklog/{workid}": (*WorkController).DeleteRequest,
			"PATCH /worklog/{workid}":  (*WorkController).PatchRequest,
		},
	}

	for path, handler := range wc.controllers.GetPaths() {
		server.HandleFunc(path, handler(wc, ctx))
	}

	wc.server = server
	return wc
}