package controllers

import (
	"context"
	"encoding/json"
	"log"
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
		log.Println("Creating work log")
		var t types.CreateWorkRequest

		json.NewDecoder(r.Body).Decode(&t)

		startDate, err := time.Parse("2006-01-02", t.Date)

		if err != nil {
			http.Error(w, "Invalid date format", http.StatusBadRequest)
			return
		}

		workID, err := wc.workService.CreateWorkLog(ctx, t.Description, startDate)
		if err != nil {
			log.Fatalf("Error starting work: %v", err)
		}

		w.Header().Set("Location", "/api/worklog/"+strconv.Itoa(workID))
		w.WriteHeader(http.StatusCreated)
	}

}

func (wc *WorkController) GetRequestById(ctx context.Context) func(http.ResponseWriter, *http.Request) {
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

		log.Println("Getting all work logs")
		// auth := r.Header.Get("Authorization")

		// if auth == "" {
		// 	http.Error(w, "Authorization header is required", http.StatusBadRequest)
		// 	return
		// }

		// stringSplit := strings.Split(auth, " ")
		// if len(stringSplit) != 2 {
		// 	http.Error(w, "Authorization header must be in the format Basic <token>", http.StatusBadRequest)
		// 	return
		// }

		// token := stringSplit[1]

		// base64Token, err := base64.StdEncoding.DecodeString(token)
		// if err != nil {
		// 	http.Error(w, "Error decoding token", http.StatusBadRequest)
		// 	return
		// }

		// tokenSplit := strings.Split(string(base64Token), ":")

		// if len(tokenSplit) != 2 {
		// 	http.Error(w, "Token must be in the format <userId>:<password>", http.StatusBadRequest)
		// 	return
		// }

		// userId := tokenSplit[0]

		// if userId == "" {
		// 	http.Error(w, "userId is required", http.StatusBadRequest)
		// 	return
		// }

		// TODO: Implement user id

		workLogs, err := wc.workService.GetAllWorkLog(ctx, 0)
		if err != nil {
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
			http.Error(w, "workId is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(workId)
		if err != nil {
			http.Error(w, "workId must be an integer", http.StatusBadRequest)
			return
		}

		log.Printf("Deleting work log by id %d", id)

		err = wc.workService.DeleteWorkLog(ctx, id)
		if err != nil {
			http.Error(w, "Error deleting work", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
func NewWorkController(ctx context.Context, server *http.ServeMux,
	workService services.WorkService,
	middleware []func(http.Handler) http.Handler) *WorkController {

	wc := &WorkController{
		workService: workService,
		controllers: ControllerPaths{
			"POST /api/worklog":            (*WorkController).PostRequest,
			"GET /api/worklog/{workid}":    (*WorkController).GetRequestById,
			"GET /api/worklog/":            (*WorkController).GetRequestAll,
			"DELETE /api/worklog/{workid}": (*WorkController).DeleteRequest,
		},
	}

	for path, handler := range wc.controllers.GetPaths() {
		server.HandleFunc(path, handler(wc, ctx))
	}

	wc.server = server
	return wc
}

func (wc *WorkController) Start() error {
	log.Printf("Starting work controller")
	return nil
}

func (wc *WorkController) Stop() error {
	log.Printf("Stopping work controller")
	return nil
}
