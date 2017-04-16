package wunderlistutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const WUNDERLIST_URL = "https://a.wunderlist.com/api/v1/"

const LIST_ENDPOINT = "lists"
const TASK_ENDPOINT = "tasks"

const LIST_PARAM = "?list_id="
const DONE_FLAG_PARAM = "&completed="

type Param struct {
	AccessToken string
	ClientID    string
}

type List struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	CreateAt string `json:"create_at"`
	ListType string `json:"list_type"`
	Revision int    `json:"revision"`
}

type Task struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	CreatedAt  string `json:"created_at"`
	DueDate    string `json:"due_date"`
	CompleteAt string `json:"completed_at"`
}

func GetLists(param Param) []List {

	// create http request
	requestString := WUNDERLIST_URL + LIST_ENDPOINT
	req, _ := http.NewRequest("GET", requestString, nil)
	req.Header.Set("X-Access-Token", param.AccessToken)
	req.Header.Set("X-Client-ID", param.ClientID)

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
		fmt.Println("unmarshal problem occured : ", err)
	}

	return lists
}

func GetTasks(param Param, listID string, doneFlag bool) []Task {

	// create http request
	requestString := WUNDERLIST_URL + TASK_ENDPOINT +
		LIST_PARAM + listID +
		DONE_FLAG_PARAM + strconv.FormatBool(doneFlag)
	req, _ := http.NewRequest("GET", requestString, nil)
	req.Header.Set("X-Access-Token", param.AccessToken)
	req.Header.Set("X-Client-ID", param.ClientID)

	// send http GET request to Wunderlist
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)

	//parse JSON
	var tasks []Task
	err = json.Unmarshal(byteArray, &tasks)
	if err != nil {
		fmt.Println("unmarshal problem occured : ", err)
	}

	return tasks
}
