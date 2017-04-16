package main

import (
	"os"
	"slackutil"
	"sort"
	"strconv"
	"strings"
	"time"

	"../wunderlistutil"
)

//for title Sort
type ByTitle []wunderlistutil.Task

func (a ByTitle) Len() int           { return len(a) }
func (a ByTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTitle) Less(i, j int) bool { return a[i].Title < a[j].Title }

func main() {

	var param wunderlistutil.Param
	param.AccessToken = os.Getenv("WUND_ACTOKEN")
	param.ClientID = os.Getenv("WUND_CLIENT")

	var listLists []wunderlistutil.List
	listLists = wunderlistutil.GetLists(param)

	doneFlag := true

	//for date
	now := time.Now()
	today := now.Format("2006-01-02")
	message := "Today's done task...\n"
	count := 0

	for _, list := range listLists {
		var taskLists []wunderlistutil.Task
		taskLists = wunderlistutil.GetTasks(param, strconv.Itoa(list.ID), doneFlag)
		sort.Sort(ByTitle(taskLists))
		applicableCount := 0
		taskString := ""
		for _, task := range taskLists {
			if strings.HasPrefix(task.CompleteAt, today) {
				count++
				applicableCount++
				taskString += "\t" + task.Title + "\n"
			}
		}
		if applicableCount == 0 {
			continue
		} else {
			message += "\t[" + list.Title + "]\n" + taskString
		}
	}

	if count == 0 {
		message += "... nothing!\nAre you seriously...?"
	} else if count < 5 {
		message += "Done tasks are few... Do it properly"
	} else {
		message += "Good job!! You done a lot of tasks!!"
	}

	slackutil.SendMessage(message)
}
