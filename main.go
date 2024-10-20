package main

import (
	"context"
	"log"
	"net/http"
	"os"
)

func startWebServer(port string, ws WorkService) error {

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

	log.Printf("Starting Work Log server on port %s\n", port)

	http.ListenAndServe(port, nil)

	return nil
}
func main() {

	port := ":3000"
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}

	workLogRepo := NewWorkLogRepository()
	workService := NewWorkService(workLogRepo)

	startWebServer(port, workService)

}
