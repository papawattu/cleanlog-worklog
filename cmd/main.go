package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	common "github.com/papawattu/cleanlog-common"
	"github.com/papawattu/cleanlog-worklog/internal/controllers"

	"github.com/papawattu/cleanlog-worklog/internal/models"
	"github.com/papawattu/cleanlog-worklog/internal/repo"
	"github.com/papawattu/cleanlog-worklog/internal/services"
)

type Config struct {
	port        string `envconfig:"PORT" default:"3000"`
	eventStore  string `envconfig:"EVENT_STORE"`
	eventStream string `envconfig:"EVENT_STREAM"`
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

	cfg := Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx := context.Background()

	var (
		workService services.WorkService
		port        string
	)

	if cfg.eventStore == "" || cfg.eventStream == "" {
		workService = services.NewWorkService(ctx, repo.NewWorkLogRepository())
	} else {
		t := common.NewHttpTransport(cfg.eventStore, cfg.eventStream, 0)
		repo := common.NewMemcacheRepository[*models.WorkLog]("localhost:11211")
		es := common.NewEventService(repo, t, "WorkLog")

		workService = services.NewWorkService(ctx, es)

		es.StartEventRunner(ctx)
	}
	if err := startWebServer(port, workService); err != nil {
		log.Fatal(err)
	}

}
