package main

import (
	"testing"
	"time"
)

func TestStartWork(t *testing.T) {
	description := "Test work log"
	wl, err := NewWorkLog(description)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if wl.WorkLogDescription != description {
		t.Errorf("Expected description %v, got %v", description, wl.WorkLogDescription)
	}
	if wl.WorkLogID != nil {
		t.Errorf("Expected WorkLogID to be nil, got %v", wl.WorkLogID)
	}
	if wl.WorkLogTimeInSecs != 0 {
		t.Errorf("Expected WorkLogTime to be 0, got %v", wl.WorkLogTimeInSecs)
	}
}

func TestLogWork(t *testing.T) {
	wl, _ := NewWorkLog("Test work log")
	task := Task{TaskID: 1}
	err := wl.LogWork(task)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(wl.Tasks) != 1 {
		t.Errorf("Expected 1 task, got %v", len(wl.Tasks))
	}
	if wl.Tasks[0].TaskID != task.TaskID {
		t.Errorf("Expected TaskID %v, got %v", task.TaskID, wl.Tasks[0].TaskID)
	}
}

func TestEndWork(t *testing.T) {
	wl, _ := NewWorkLog("Test work log")
	time.Sleep(1 * time.Second) // Simulate some work time
	err := wl.EndWork()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if wl.WorkLogTimeInSecs == 0 {
		t.Errorf("Expected WorkLogTime to be greater than 0, got %v", wl.WorkLogTimeInSecs)
	}
}
