package models

import (
	"testing"
	"time"
)

func TestStartWork(t *testing.T) {
	description := "Test work log"
	wl, err := NewWorkLog(description, time.Now())
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
	wl, _ := NewWorkLog("Test work log", time.Now())
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
	wl, _ := NewWorkLog("Test work log", time.Now())
	time.Sleep(1 * time.Second) // Simulate some work time
	err := wl.EndWork()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if wl.WorkLogTimeInSecs == 0 {
		t.Errorf("Expected WorkLogTime to be greater than 0, got %v", wl.WorkLogTimeInSecs)
	}
}

func TestMultipleTasks(t *testing.T) {
	wl, _ := NewWorkLog("Test work log", time.Now())
	tasks := []Task{
		{TaskID: 1},
		{TaskID: 2},
		{TaskID: 3},
	}
	for _, task := range tasks {
		err := wl.LogWork(task)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	}
	if len(wl.Tasks) != len(tasks) {
		t.Errorf("Expected %v tasks, got %v", len(tasks), len(wl.Tasks))
	}
	for i, task := range tasks {
		if wl.Tasks[i].TaskID != task.TaskID {
			t.Errorf("Expected TaskID %v, got %v", task.TaskID, wl.Tasks[i].TaskID)
		}
	}
}

func TestEndWorkWithoutLoggingTasks(t *testing.T) {
	wl, _ := NewWorkLog("Test work log", time.Now())
	time.Sleep(1 * time.Second) // Simulate some work time
	err := wl.EndWork()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if wl.WorkLogTimeInSecs == 0 {
		t.Errorf("Expected WorkLogTime to be greater than 0, got %v", wl.WorkLogTimeInSecs)
	}
	if len(wl.Tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %v", len(wl.Tasks))
	}
}
func TestAddTask(t *testing.T) {
	wl, _ := NewWorkLog("Test work log", time.Now())
	task := Task{TaskID: 1}
	err := wl.AddTask(task)
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

func TestRemoveTask(t *testing.T) {
	wl, _ := NewWorkLog("Test work log", time.Now())
	task := Task{TaskID: 1}
	err := wl.AddTask(task)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	err = wl.RemoveTask(task)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(wl.Tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %v", len(wl.Tasks))
	}
}

func TestChangeDescription(t *testing.T) {
	wl, _ := NewWorkLog("Test work log", time.Now())
	description := "New description"
	err := wl.ChangeDescription(description)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if wl.WorkLogDescription != description {
		t.Errorf("Expected description %v, got %v", description, wl.WorkLogDescription)
	}
}

func TestTasksAreInitialized(t *testing.T) {
	wl, _ := NewWorkLog("Test work log", time.Now())
	if wl.Tasks == nil {
		t.Errorf("Expected Tasks to be initialized, got nil")
	}
}
