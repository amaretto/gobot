package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nlopes/slack"
)

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

//for title Sort
type ByTitle []Task

func (a ByTitle) Len() int           { return len(a) }
func (a ByTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTitle) Less(i, j int) bool { return a[i].Title < a[j].Title }

func main() {

	// Get information for connecting Wunderlist
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
	var listID string
	doneFlagParam := "&completed="
	doneFlag := "true"

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
	var tasks []Task
	err = json.Unmarshal(taskByteArray, &tasks)
	if err != nil {
		fmt.Println("err! : ", err)
	}

	//sort by title=
	sort.Sort(ByTitle(tasks))
	//for date
	now := time.Now()
	today := now.Format("2006-01-02")
	message := "Today's done task...\n"
	count := 0

	for _, task := range tasks {
		if strings.HasPrefix(task.CompleteAt, today) {
			count++
			fmt.Println("\t", task.Title)
			message += task.Title + "\n"
		}
	}

	if count == 0 {
		message += "... nothing!\nAre you seriously...?"
	} else if count < 5 {
		message += "Done tasks are few... Do it properly"
	} else {
		message += "Good job!! You done a lot of tasks!!"
	}

	sendMessage(message)
}

func sendMessage(message string) {
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	username := os.Getenv("SLACK_USERNAME")

	message = username + " " + message
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Text: message,
	}
	params.Attachments = []slack.Attachment{attachment}
	params.AsUser = true
	api.PostMessage("#bot_project", "", params)
}
