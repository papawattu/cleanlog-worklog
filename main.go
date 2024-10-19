package main

import (
	"context"
	"log"
	"net/http"
)

func startWebServer(ws WorkService) error {

	controllers, err := MakeControllers(context.Background(), NewWorkController(ws))
	if err != nil {
		log.Fatal(err)
		return err
	}

	if err = controllers.Start(); err != nil {
		log.Fatal(err)
		return err
	}

	http.HandleFunc("/api/worklog/{workid}", controllers.HandleRequest)
	http.HandleFunc("/api/worklog", controllers.HandleRequest)

	http.ListenAndServe(":3000", nil)

	return nil
}
func main() {

	workLogRepo := NewWorkLogRepository()
	workService := NewWorkService(workLogRepo)

	startWebServer(workService)

	workID, err := workService.StartWork("Test work log")
	if err != nil {
		log.Fatalf("Error starting work: %v", err)
	}

	task := Task{TaskID: 1}
	err = workService.LogWork(workID, task)
	if err != nil {
		log.Fatalf("Error logging work: %v", err)
	}

	err = workService.EndWork(workID)
	if err != nil {
		log.Fatalf("Error ending work: %v", err)
	}

	wl, err := workService.GetWorkLog(workID)

	if err != nil {
		log.Fatalf("Error getting work log: %v", err)
	}

	log.Printf("Work log Id: %+v", *wl.WorkLogID)

}
