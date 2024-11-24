package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	common "github.com/papawattu/cleanlog-common"
	"github.com/papawattu/cleanlog-worklog/internal/models"
	"github.com/papawattu/cleanlog-worklog/internal/services"
	"github.com/papawattu/cleanlog-worklog/types"
)

var (
	location     string
	taskLocation string
	workLogRepo  = common.NewInMemoryRepository[*models.WorkLog]()
	workService  = services.NewWorkService(context.Background(), workLogRepo)
)

func GetWorkLogByID(id int) *models.WorkLog {
	wl, err := workService.GetWorkLog(context.Background(), id)
	if err != nil {
		return nil
	}
	return wl
}
func TestGetController(t *testing.T) {

	ctx := context.Background()

	ctx = context.WithValue(ctx, "user", 0)

	controllers := NewWorkController(ctx, http.NewServeMux(), workService)

	server := httptest.NewServer(controllers.server)

	defer server.Close()

	r, err := http.Get(server.URL + "/api/worklog/1")

	if err != nil {
		t.Fatal(err)
	}

	if r.StatusCode != http.StatusNotFound {
		t.Fatalf(t.Name()+"Expected status code %v, got %v", http.StatusNotFound, r.StatusCode)
	}

	t.Log("Test passed")

}

func TestCreateWorkLogController(t *testing.T) {

	ctx := context.Background()

	ctx = context.WithValue(ctx, "user", 0)

	controllers := NewWorkController(ctx, http.NewServeMux(), workService)

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

	ctx := context.Background()

	ctx = context.WithValue(ctx, "user", 0)

	controllers := NewWorkController(ctx, http.NewServeMux(), workService)

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

func TestUpdateWorkLogController(t *testing.T) {

	if location == "" {
		t.Error("Location is not set")
		t.FailNow()
		return
	}

	ctx := context.Background()

	ctx = context.WithValue(ctx, "user", 0)

	controllers := NewWorkController(ctx, http.NewServeMux(), workService)

	server := httptest.NewServer(controllers.server)

	defer server.Close()

	jsonStr := []byte(`{"description":"Updated work log", "date":"2021-01-01"}`)

	req, err := http.NewRequest(http.MethodPut, server.URL+location, bytes.NewReader(jsonStr))

	if err != nil {
		t.Fatal(err)
	}

	r, err := server.Client().Do(req)

	if err != nil {
		t.Fatal(err)
	}

	if r.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status code %v, got %v", http.StatusNoContent, r.StatusCode)
	}

}
func TestAddTaskToWorkLogController(t *testing.T) {

	if location == "" {
		t.Error("Location is not set")
		t.FailNow()
		return
	}

	ctx := context.Background()

	ctx = context.WithValue(ctx, "user", 0)

	controllers := NewWorkController(ctx, http.NewServeMux(), workService)

	server := httptest.NewServer(controllers.server)

	defer server.Close()

	jsonStr := []byte(`{"taskId":123}`)

	t.Logf("Location: %s", server.URL+location)

	s := strings.Split(location, "/")
	id := s[len(s)-1]
	idInt, err := strconv.Atoi(id)

	if err != nil {
		t.Fatal(err)
	}
	wl := GetWorkLogByID(idInt)

	if wl == nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, server.URL+location+"/task", bytes.NewReader(jsonStr))

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

	taskLocation = r.Header.Get("Location")

	t.Logf("Task Location: %s", location)

	if taskLocation == "" {
		t.Fatalf("Expected Location header to be set")
	}
}
func TestDeleteTaskFromWorkLogController(t *testing.T) {

	if location == "" || taskLocation == "" {
		t.Error("Locations are not set")
		t.FailNow()
		return
	}
	ctx := context.Background()

	ctx = context.WithValue(ctx, "user", 0)

	controllers := NewWorkController(ctx, http.NewServeMux(), workService)

	server := httptest.NewServer(controllers.server)

	defer server.Close()

	req, err := http.NewRequest(http.MethodDelete, server.URL+taskLocation, nil)

	if err != nil {
		t.Fatal(err)
	}

	r, err := server.Client().Do(req)

	if err != nil {
		t.Fatal(err)
	}

	if r.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status code %v, got %v", http.StatusNoContent, r.StatusCode)
	}

}
func TestDeleteWorkLogController(t *testing.T) {

	if location == "" {
		t.Error("Location is not set")
		t.FailNow()
		return
	}
	ctx := context.Background()

	ctx = context.WithValue(ctx, "user", 0)

	controllers := NewWorkController(ctx, http.NewServeMux(), workService)

	server := httptest.NewServer(controllers.server)

	defer server.Close()

	req, err := http.NewRequest(http.MethodDelete, server.URL+location, nil)

	if err != nil {
		t.Fatal(err)
	}

	r, err := server.Client().Do(req)

	if err != nil {
		t.Fatal(err)
	}

	if r.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status code %v, got %v", http.StatusNoContent, r.StatusCode)
	}

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
