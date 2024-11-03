package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

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
	req.Header.Set("Authorization", "Authorization: Basic amFtaWU6c2ltcHNvbnM=")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal("Error:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Fatalf("Error: status code %d\n", resp.StatusCode)
		return "", nil
	}

	log.Printf("Work log created with ID: %s\n", resp.Header.Get("Location"))

	return resp.Header.Get("Location"), nil
}
func GetWorkLog(loc string, baseUri string) {
	url := fmt.Sprintf("%s%s", baseUri, loc)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error: status code", resp.StatusCode)
		return
	}

	r := &types.WorkResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		log.Println("Error decoding JSON:", err)
		return
	}

	log.Printf("Work log: %v", r)

}

func GetAllWorkLogs(baseUri string) {
	url := fmt.Sprintf("%s/api/worklog/", baseUri)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Authorization: Basic amFtaWU6c2ltcHNvbnM=")

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

}
