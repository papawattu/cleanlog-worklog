package events

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/papawattu/cleanlog-worklog/internal/models"
	"github.com/papawattu/cleanlog-worklog/internal/repo"
)

const (
	WorkLogCreated = "WorkLogCreated"
	WorkLogDeleted = "WorkLogDeleted"
	EventUri       = "/event"
	EventVersion   = 1
)

type EventBroadcaster struct {
	repo         repo.WorkLogRepository
	broadcastUri string
}

type Event struct {
	EventType    string    `json:"eventType"`
	EventTime    time.Time `json:"eventTime"`
	EventVersion uint32    `json:"eventVersion"`
	EventSHA     string    `json:"eventSHA"`
	EventData    string    `json:"eventData"`
}

func (eb *EventBroadcaster) postEvent(event Event) error {

	ev, err := json.Marshal(event)
	if err != nil {
		return err
	}

	r, err := http.Post(eb.broadcastUri, "application/json", bytes.NewBuffer(ev))

	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusCreated {
		return fmt.Errorf("Error: status code %d", r.StatusCode)
	}

	return nil
}

func (eb *EventBroadcaster) SaveWorkLog(wl *models.WorkLog) error {

	wlj, err := json.Marshal(wl)

	if err != nil {
		return err
	}

	h := sha256.New()

	h.Write([]byte(wlj))

	// Broadcast event
	event := Event{
		EventType:    WorkLogCreated,
		EventTime:    time.Now(),
		EventVersion: EventVersion,
		EventSHA:     fmt.Sprintf("%x", h.Sum(nil)),
		EventData:    string(wlj),
	}

	err = eb.postEvent(event)

	if err != nil {
		return err
	}

	err = eb.repo.SaveWorkLog(wl)

	if err != nil {

		err := eb.DeleteWorkLog(*wl.WorkLogID)
		if err != nil {
			log.Panicf("Error saving work log: %v published rollback event", err)
		}
		log.Printf("Error saving work log: %v published rollback event", err)
		return err
	}

	return nil //
}

func (eb *EventBroadcaster) GetWorkLog(id int) (*models.WorkLog, error) {
	return eb.repo.GetWorkLog(id)
}

func (eb *EventBroadcaster) GetAllWorkLogsForUser(userID int) ([]*models.WorkLog, error) {
	return eb.repo.GetAllWorkLogsForUser(userID)
}

func (eb *EventBroadcaster) DeleteWorkLog(id int) error {

	wl, err := eb.repo.GetWorkLog(id)

	if err != nil {
		return err
	}

	wlj, err := json.Marshal(wl)

	if err != nil {
		return err
	}

	h := sha256.New()

	h.Write([]byte(wlj))

	// Broadcast event
	event := Event{
		EventType:    WorkLogDeleted,
		EventTime:    time.Now(),
		EventVersion: EventVersion,
		EventSHA:     fmt.Sprintf("%x", h.Sum(nil)),
		EventData:    string(wlj),
	}

	err = eb.postEvent(event)

	if err != nil {
		return err
	}

	return eb.repo.DeleteWorkLog(id)
}

func NewEventBroadcaster(repo repo.WorkLogRepository, baseUri string, topic string) *EventBroadcaster {
	return &EventBroadcaster{
		repo:         repo,
		broadcastUri: baseUri + "/event/" + topic,
	}
}
