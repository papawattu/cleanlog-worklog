package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/papawattu/cleanlog-worklog/internal/controllers"
	"github.com/papawattu/cleanlog-worklog/internal/events"
	"github.com/papawattu/cleanlog-worklog/internal/repo"
	"github.com/papawattu/cleanlog-worklog/internal/services"
)

func startWebServer(server *http.Server, ws services.WorkService) error {

	// make a middleware array

	// middleware := []func(http.Handler) http.Handler{
	// 	func(next http.Handler) http.Handler {
	// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 			log.Printf("Request received: %s %s\n", r.Method, r.URL)
	// 			next.ServeHTTP(w, r)
	// 			log.Printf("Request completed: %s %s\n", r.Method, r.URL)
	// 		})
	// 	},
	// }

	controllers.NewWorkController(context.Background(), server.Handler.(*http.ServeMux), ws, func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s %s\n", r.Method, r.URL)
	})

	log.Printf("Starting Work Log server on port %s\n", server.Addr)

	return server.ListenAndServe()

}
func main() {

	var (
		workService services.WorkService
		port        string
		server      *http.Server
	)

	port = "3000"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	topic := "worklog"

	if os.Getenv("EVENT_TOPIC") != "" {
		topic = os.Getenv("EVENT_TOPIC")
	}

	if os.Getenv("EVENT_BROADCASTER") != "" {
		eventBroadcaster := events.NewEventBroadcaster(repo.NewWorkLogRepository(), os.Getenv("EVENT_BROADCASTER"), topic)
		workService = services.NewWorkService(eventBroadcaster)
	} else {
		workService = services.NewWorkService(repo.NewWorkLogRepository())
	}
	server = &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: http.NewServeMux(),
	}

	if err := startWebServer(server, workService); err != nil {
		log.Fatal(err)
	}

}
