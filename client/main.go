package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/papawattu/cleanlog-worklog/types"
)

func CreateWorkLog(description string, baseUri string) int {
	url := fmt.Sprintf("%s/api/worklog", baseUri)
	body := types.CreateWorkRequest{Description: description}
	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return -1
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Error:", err)
		return -1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Error: status code", resp.StatusCode)
		return 0
	}

	r := map[string]int{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return 0
	}

	fmt.Println("Work log created with ID:", r["workId"])

	return r["workId"]
}
func GetWorkLog(id int, baseUri string) {
	url := fmt.Sprintf("%s/api/worklog/%d", baseUri, id)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		return
	}

	r := &types.CreateWorkResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Println("Work log:", r)

}

func main() {
	id := CreateWorkLog("Work log 1", "http://localhost:3000")

	GetWorkLog(id, "http://localhost:3000")
}
