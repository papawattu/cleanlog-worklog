package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/papawattu/cleanlog-worklog/types"
)

type WorkController struct {
	workService WorkService
	handlePost  func(w http.ResponseWriter, r *http.Request)
	handleGet   func(w http.ResponseWriter, r *http.Request)
}

func NewWorkController(workService WorkService) *WorkController {

	return &WorkController{workService: workService, handlePost: func(w http.ResponseWriter, r *http.Request) {
		var t types.CreateWorkRequest

		json.NewDecoder(r.Body).Decode(&t)

		workID, err := workService.StartWork(t.Description)
		if err != nil {
			log.Fatalf("Error starting work: %v", err)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"workId": workID})
	}, handleGet: func(w http.ResponseWriter, r *http.Request) {
		workId := r.PathValue("workid")
		if workId == "" {
			http.Error(w, "workId is required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(workId)
		if err != nil {
			http.Error(w, "workId must be an integer %s", http.StatusBadRequest)
			return
		}

		wl, err := workService.GetWorkLog(id)
		if err != nil {
			log.Fatalf("Error getting work log: %v", err)
			http.Error(w, "Error getting work log", http.StatusInternalServerError)
		}

		if wl == nil {
			http.Error(w, "Work log not found", http.StatusNotFound)
		} else {
			taskIds := []int{}

			for _, t := range wl.Tasks {
				taskIds = append(taskIds, t.TaskID)
			}

			r := types.CreateWorkResponse{WorkID: *wl.WorkLogID, Description: wl.WorkLogDescription, TaskIds: taskIds}
			json.NewEncoder(w).Encode(r)
		}

	}}
}

func (wc *WorkController) Start() error {
	log.Printf("Starting work controller")
	return nil
}

func (wc *WorkController) Stop() error {
	log.Printf("Stopping work controller")
	return nil
}

func (wc *WorkController) HandleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling request: %v", r)

	switch {
	case r.Method == "POST":
		wc.handlePost(w, r)
	case r.Method == "GET":
		wc.handleGet(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}
