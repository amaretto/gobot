package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type List struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	CreateAt string `json:"create_at"`
	ListType string `json:"list_type"`
	Revision int    `json:"revision"`
}

/*  */
type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	DueDate   string `json:"due_date"`
	CompeteAt string `json:"completed_at"`
}

func main() {

	// Get information for connecting Wunderlist from env
	accessToken := os.Getenv("WUND_ACTOKEN")
	clientID := os.Getenv("WUND_CLIENT")

	wunderlistURL := "https://a.wunderlist.com/api/v1/"

	// get all lists from Wunderlist
	endpoint := "lists"

	// create http request
	req, _ := http.NewRequest("GET", wunderlistURL+endpoint, nil)
	req.Header.Set("X-Access-Token", accessToken)
	req.Header.Set("X-Client-ID", clientID)

	// send http GET request to Wunderlist
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)

	//parse JSON
	var lists []List
	err = json.Unmarshal(byteArray, &lists)
	if err != nil {
		fmt.Println("Unmarshal Problem? : ", err)
	}

	taskEndpoint := "tasks"
	listParam := "?list_id="
	doneFlagParam := "&completed="
	doneFlag := "true"

	for _, list := range lists {
		fmt.Println(list.Title)

		if list.Title == "job" {
			// create http request
			taskReq, _ := http.NewRequest("GET", wunderlistURL+
				taskEndpoint+
				listParam+
				strconv.Itoa(list.ID)+
				doneFlagParam+doneFlag, nil)
			taskReq.Header.Set("X-Access-Token", accessToken)
			taskReq.Header.Set("X-Client-ID", clientID)
			fmt.Println(taskReq)
			//get task list from Wunderlist
			taskResp, err := client.Do(taskReq)
			defer taskResp.Body.Close()

			taskByteArray, _ := ioutil.ReadAll(taskResp.Body)

			//parse task list JSON
			var tasks []Task
			err = json.Unmarshal(taskByteArray, &tasks)
			if err != nil {
				fmt.Println("err! : ", err)
			}

			for _, task := range tasks {
				fmt.Println("\t", task.Title)
			}
		}

	}
}
