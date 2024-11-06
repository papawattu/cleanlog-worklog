package events

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"strings"

	utils "github.com/papawattu/cleanlog-worklog/internal"
)

func EventStream(ctx context.Context, baseUri string, es chan string, topic string) {
	for {
		client := utils.NewRetryableClient(10)

		req, err := http.NewRequest("GET", baseUri+"/eventstream/"+topic, nil)

		if err != nil {
			log.Fatalf("Error creating request: %v", err)
		}

		req.Header.Set("Accept", "text/event-stream")
		req.Header.Set("Cache-Control", "no-cache")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error connecting to event stream: %v", err)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Error: status code %d", resp.StatusCode)
			es <- "Error connecting to event stream"
		}

		log.Println("Connected to event stream")
		scanner := bufio.NewScanner(resp.Body)

		for running := true; running; {
			select {
			case <-ctx.Done():
				log.Println("Timeout")
				running = false
				break
			case <-req.Context().Done():
				log.Println("Client connection closed")
				running = false
				break
			default:
				scanner.Scan()
				e := scanner.Text()
				if e == "" {
					log.Println("Empty event")
					running = false
					break
				}

				if !strings.HasPrefix(e, "data: ") {
					log.Fatalf("Error: unexpected event %s", e)
					running = false
					break
				}
				es <- strings.TrimLeft(e, "data: ")
			}
		}
	}

}
