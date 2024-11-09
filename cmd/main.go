package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/papawattu/cleanlog-worklog/internal/controllers"
	"github.com/papawattu/cleanlog-worklog/internal/events"
	"github.com/papawattu/cleanlog-worklog/internal/middleware"
	"github.com/papawattu/cleanlog-worklog/internal/repo"
	"github.com/papawattu/cleanlog-worklog/internal/services"
)

func startWebServer(port string, ws services.WorkService) error {
	stack := middleware.CreateMiddleware(middleware.Logging)

	router := http.NewServeMux()

	api := http.NewServeMux()
	api.Handle("/api/", http.StripPrefix("/api", router))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: stack(api),
	}

	controllers.NewWorkController(context.Background(), router, ws)

	log.Printf("Starting Work Log server on port %s\n", port)
	return server.ListenAndServe()

}
func main() {

	ctx := context.Background()

	var (
		workService services.WorkService
		port        string
	)

	port = "3000"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	topic := "worklog"

	if os.Getenv("EVENT_TOPIC") != "" {
		topic = os.Getenv("EVENT_TOPIC")
	}

	if os.Getenv("EVENT_BROADCASTER") == "" || os.Getenv("EVENT_STREAM") == "" {
		workService = services.NewWorkService(ctx, repo.NewWorkLogRepository())
	} else {
		eventBroadcaster := events.NewEventBroadcaster(ctx, repo.NewWorkLogRepository(), os.Getenv("EVENT_BROADCASTER"), os.Getenv("EVENT_STREAM"), topic)
		workService = services.NewWorkService(ctx, eventBroadcaster)
	}

	if err := startWebServer(port, workService); err != nil {
		log.Fatal(err)
	}

}
