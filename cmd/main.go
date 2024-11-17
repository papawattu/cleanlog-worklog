package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	common "github.com/papawattu/cleanlog-common"
	"github.com/papawattu/cleanlog-worklog/internal/controllers"
	"github.com/papawattu/cleanlog-worklog/internal/middleware"
	"github.com/papawattu/cleanlog-worklog/internal/repo"
	"github.com/papawattu/cleanlog-worklog/internal/services"
)

func startWebServer(port string, ws services.WorkService) error {
	stack := middleware.CreateMiddleware(
		middleware.Recover,
		middleware.Logging,
		middleware.Authenticated,
	)

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

	if os.Getenv("EVENT_STORE") == "" || os.Getenv("EVENT_STREAM") == "" {
		workService = services.NewWorkService(ctx, repo.NewWorkLogRepository())
	} else {
		t := common.NewHttpTransport(os.Getenv("EVENT_STORE"), os.Getenv("EVENT_STREAM"), 0)

		es := common.NewEventService(repo.NewWorkLogRepository(), t, "WorkLog")

		workService = services.NewWorkService(ctx, es)

		//es.StartEventRunner(ctx)

		go func() {

			slog.Info("Starting event runner")

			err := es.Connect(ctx)

			if err != nil {
				slog.Error("Error connecting to event store", "error", err)
				return
			}
			for {
				select {
				case <-ctx.Done():
					return
				default:
					slog.Info("Waiting for event")
					ev, err := es.NextEvent()

					slog.Info("Got event", "event", ev, "error", err)
					if err != nil {
						log.Fatal(err)
					}

					if ev != nil {
						slog.Info("Handling event", "event", *ev)
						es.HandleEvent(*ev)
					} else {
						slog.Info("No event")
					}
				}
			}
		}()
	}

	if err := startWebServer(port, workService); err != nil {
		log.Fatal(err)
	}

}
