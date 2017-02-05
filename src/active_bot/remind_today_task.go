package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"../slackutil"
	"../wunderlistutil"
)

//for title Sort
type ByTitle []wunderlistutil.Task

func (a ByTitle) Len() int           { return len(a) }
func (a ByTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTitle) Less(i, j int) bool { return a[i].Title < a[j].Title }

//for avoiding invalid due date
const RESETED_DUE_DATE = "0001-01-01"

func main() {

	var param wunderlistutil.Param
	param.AccessToken = os.Getenv("WUND_ACTOKEN")
	param.ClientID = os.Getenv("WUND_CLIENT")

	// Get information for connecting Wunderlist
	accessToken := os.Getenv("WUND_ACTOKEN")
	clientID := os.Getenv("WUND_CLIENT")
	wunderlistURL := "https://a.wunderlist.com/api/v1/"
	//
	//	// get all lists from Wunderlist
	//	endpoint := "lists"
	//
	//	// create http request
	//	req, _ := http.NewRequest("GET", wunderlistURL+endpoint, nil)
	//	req.Header.Set("X-Access-Token", accessToken)
	//	req.Header.Set("X-Client-ID", clientID)
	//
	//	// send http GET request to Wunderlist
	client := new(http.Client)
	//	resp, err := client.Do(req)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	defer resp.Body.Close()
	//
	//	byteArray, _ := ioutil.ReadAll(resp.Body)
	//
	//	//parse JSON
	var lists []wunderlistutil.List
	//	err = json.Unmarshal(byteArray, &lists)
	//	if err != nil {
	//		fmt.Println("Unmarshal Problem? : ", err)
	//	}
	lists = wunderlistutil.GetLists(param)

	taskEndpoint := "tasks"
	listParam := "?list_id="
	var listID string
	doneFlagParam := "&completed="
	doneFlag := "false"

	for _, list := range lists {
		if list.Title == "job" {
			listID = strconv.Itoa(list.ID)
		}
	}

	// create http request
	taskReq, _ := http.NewRequest("GET", wunderlistURL+
		taskEndpoint+
		listParam+
		listID+
		doneFlagParam+doneFlag, nil)
	taskReq.Header.Set("X-Access-Token", accessToken)
	taskReq.Header.Set("X-Client-ID", clientID)
	//get task list from Wunderlist
	taskResp, err := client.Do(taskReq)
	defer taskResp.Body.Close()

	taskByteArray, _ := ioutil.ReadAll(taskResp.Body)

	//parse task list JSON
	var tasks []wunderlistutil.Task
	err = json.Unmarshal(taskByteArray, &tasks)
	if err != nil {
		fmt.Println("err! : ", err)
	}

	//sort by title=
	sort.Sort(ByTitle(tasks))

	//for date
	now := time.Now()
	layout := "2006-01-02"
	todayString := now.Format(layout)
	today, err := time.Parse(layout, todayString)
	headMessage := "Today's your task...\n"
	newTaskMessage := "\t[New Task]\n"
	delayedTaskMessage := "\t[Delayed Task]\n"
	var dueDate time.Time
	var dueDateString string

	//make message string
	for _, task := range tasks {
		dueDate, err = time.Parse(layout, task.DueDate)
		dueDateString = dueDate.Format(layout)
		if dueDate.Equal(today) {
			newTaskMessage += "\t" + task.Title + "\n"
		} else if dueDate.Before(today) && (dueDateString != RESETED_DUE_DATE) {
			delayedTaskMessage += "\t" + task.Title + "(" + dueDateString + ")\n"
		}
	}

	slackutil.SendMessage(headMessage + newTaskMessage + delayedTaskMessage)
}
