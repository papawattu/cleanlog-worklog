package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	common "github.com/papawattu/cleanlog-common"
	"github.com/papawattu/cleanlog-worklog/internal/controllers"
	"github.com/papawattu/cleanlog-worklog/internal/models"

	"github.com/papawattu/cleanlog-worklog/internal/services"
)

type Config struct {
	Port        string `envconfig:"PORT" default:"3000"`
	EventStore  string `envconfig:"EVENT_STORE"`
	EventStream string `envconfig:"EVENT_STREAM"`
}

func startWebServer(port string, ws services.WorkService) error {

	stack := common.CreateMiddleware(
		common.Recover,
		common.Logging,
		common.Authenticated,
	)

	router := http.NewServeMux()

	api := http.NewServeMux()
	api.Handle("/", router)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: stack(api),
	}

	controllers.NewWorkController(context.Background(), router, ws)

	log.Printf("Starting Work Log server on port %s\n", port)
	return server.ListenAndServe()

}
func main() {

	var cfg Config

	err := envconfig.Process("worklog", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx := context.Background()

	var (
		workService services.WorkService
	)

	if cfg.EventStore == "" || cfg.EventStream == "" {
		workService = services.NewWorkService(ctx, common.NewInMemoryRepository[*models.WorkLog]())
	} else {
		t := common.NewHttpTransport(cfg.EventStore, cfg.EventStream, 0)
		repo := common.NewMemcacheRepository[*models.WorkLog]("localhost:11211", "worklog", nil)
		es := common.NewEventService(repo, t, "WorkLog")

		workService = services.NewWorkService(ctx, es)

		es.StartEventRunner(ctx)
	}
	slog.Info("Starting Work Log server", "port", cfg.Port)
	if err := startWebServer(cfg.Port, workService); err != nil {
		log.Fatal(err)
	}

}
