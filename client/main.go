package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/papawattu/cleanlog-worklog/types"
)

func CreateWorkLog(description string, baseUri string) (string, error) {
	url := fmt.Sprintf("%s/api/worklog", baseUri)
	body := types.CreateWorkRequest{Description: description, Date: "2024-01-01"}
	b, err := json.Marshal(body)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		log.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mytoken")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal("Error:", err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusCreated {
		log.Fatalf("Error: status code %d\n", resp.StatusCode)
		return "", nil
	}

	log.Printf("Work log created with ID: %s\n", resp.Header.Get("Location"))

	return resp.Header.Get("Location"), nil
}
func GetWorkLog(loc string, baseUri string) {
	url := fmt.Sprintf("%s%s", baseUri, loc)

	var count int = 0

	for {

		req, err := http.NewRequest(http.MethodGet, url, nil)

		if err != nil {
			log.Println("Error:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer mytoken")

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			log.Println("Error:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			if count > 20 {
				log.Fatalln("Work log not found")
				return
			}
			count++
			log.Printf("Work log not found at %s waiting - times %d\n", url, count)
			time.Sleep(1 * time.Second)
		} else {
			if resp.StatusCode == http.StatusOK {
				log.Printf("Work log found at %s\n", url)
				r := &types.WorkResponse{}
				if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
					log.Println("Error decoding JSON:", err)
					return
				}

				break
			} else {
				log.Fatalf("Error: status code %d\n", resp.StatusCode)
				return
			}
		}
	}
}

func GetAllWorkLogs(baseUri string) {
	url := fmt.Sprintf("%s/api/worklog/", baseUri)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mytoken")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Println("Error: status code", resp.StatusCode)
		return
	}

	r := types.ListWorkResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Println("Error decoding JSON:", err)
		return
	}

	log.Printf("Work logs: %+v", r)
}

func DeleteWorkLog(loc string, baseUri string) {
	url := fmt.Sprintf("%s%s", baseUri, loc)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Fatalln("Error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mytoken")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("Error:", err)
		return
	}

	log.Println("Response status code:", resp.StatusCode)
	if resp.StatusCode != http.StatusNoContent {
		log.Fatalln("Error: status code ", resp.StatusCode)
		return
	}

	log.Println("Work log deleted")
}
func main() {
	var baseUri string
	flag.StringVar(&baseUri, "baseUri", "http://localhost:3000", "Base URI for the worklog service")
	flag.Parse()

	log.Println("Creating work log")
	id, err := CreateWorkLog("Work log 1", baseUri)
	if err != nil {
		log.Fatalf("Error creating work log: %v", err)
	}

	log.Printf("Getting work log %s\n", id)
	GetWorkLog(id, baseUri)

	log.Println("Getting all work logs")
	GetAllWorkLogs(baseUri)

	log.Printf("Deleting work log %s\n", id)
	DeleteWorkLog(id, baseUri)

	// ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)

	// defer cancel()

	// ch := make(chan string)

	// topic := "worklog"
	// go events.EventStream(ctx, "http://localhost:8090", ch, topic)

	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		log.Println("Timeout")
	// 		return
	// 	case e := <-ch:
	// 		log.Println(e)
	// 	default:
	// 		//

	// 	}
	// }

}
