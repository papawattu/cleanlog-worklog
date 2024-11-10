package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/papawattu/cleanlog-worklog/internal/models"
	"github.com/papawattu/cleanlog-worklog/internal/repo"
	"github.com/papawattu/cleanlog-worklog/internal/services"
	"github.com/papawattu/cleanlog-worklog/types"
)

var (
	location    string
	workLogRepo = repo.NewWorkLogRepository()
	workService = services.NewWorkService(context.Background(), workLogRepo)
)

func TestGetController(t *testing.T) {

	controllers := NewWorkController(context.Background(), http.NewServeMux(), workService)

	server := httptest.NewServer(controllers.server)

	defer server.Close()

	r, err := http.Get(server.URL + "/api/worklog/")

	if err != nil {
		t.Fatal(err)
	}

	if r.StatusCode != http.StatusOK {
		t.Fatalf(t.Name()+"Expected status code %v, got %v", http.StatusOK, r.StatusCode)
	}

	t.Log("Test passed")

}

func TestCreateWorkLogController(t *testing.T) {

	controllers := NewWorkController(context.Background(), http.NewServeMux(), workService)

	server := httptest.NewServer(controllers.server)

	defer server.Close()

	jsonStr := []byte(`{"description":"Test work log", "date":"2021-01-01"}`)

	req, err := http.NewRequest(http.MethodPost, server.URL+"/api/worklog", bytes.NewReader(jsonStr))

	if err != nil {
		t.Fatal(err)
	}

	r, err := server.Client().Do(req)

	if err != nil {
		t.Fatal(err)
	}

	if r.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status code %v, got %v", http.StatusCreated, r.StatusCode)
	}

	location = r.Header.Get("Location")

	if location == "" {
		t.Fatalf("Expected Location header to be set")
	}
}
func TestGetWorkLogController(t *testing.T) {
	if location == "" {
		t.Skip("Skipping TestGetWorkLogController because location is not set")
	}

	controllers := NewWorkController(context.Background(), http.NewServeMux(), workService)

	server := httptest.NewServer(controllers.server)

	defer server.Close()

	t.Logf("Location: %s", server.URL+location)

	r, err := http.Get(server.URL + location)

	if err != nil {
		t.Fatal(err)
	}

	if r.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %+v, got %+v", http.StatusOK, r.StatusCode)
	}

	var wlr types.WorkResponse

	err = json.NewDecoder(r.Body).Decode(&wlr)

	if err != nil {
		t.Fatal(err)
	}

	if wlr.Description != "Test work log" {
		t.Fatalf("Expected WorkLogDescription to be 'Test work log', got %v", wlr.Description)
	}

	if wlr.Date != "2021-01-01" {
		t.Fatalf("Expected WorkLogDate to be '2021-01-01', got %v", wlr.Date)
	}
	t.Log("Test passed")

}
func TestInlineTasksWithEmptyTasks(t *testing.T) {
	tasks := []models.Task{}

	taskIds := inlineTasks(tasks)

	if !reflect.DeepEqual(taskIds, []int{}) {
		t.Fatalf("Expected empty int array, got %v", taskIds)
	}
	t.Log("Test passed")
}
func TestInlineTasks(t *testing.T) {
	tasks := []models.Task{{TaskID: 1}, {TaskID: 2}, {TaskID: 3}}

	taskIds := inlineTasks(tasks)

	if !reflect.DeepEqual(taskIds, []int{1, 2, 3}) {
		t.Fatalf("Expected '1, 2, 3', got %v", taskIds)
	}
	t.Log("Test passed")
}
func TestNilTasks(t *testing.T) {

	taskIds := inlineTasks(nil)

	if !reflect.DeepEqual(taskIds, []int{}) {
		t.Fatalf("Expected empty int array, got %v", taskIds)
	}
	t.Log("Test passed")
}