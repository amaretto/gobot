package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

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

	var lists []wunderlistutil.List
	lists = wunderlistutil.GetLists(param)

	//for date
	now := time.Now()
	layout := "2006-01-02"
	todayString := now.Format(layout)
	today, err := time.Parse(layout, todayString)
	if err != nil {
		fmt.Println("err! : ", err)
	}
	headMessage := "Today's your task...\n"
	newTaskMessage := "\t\t[New]"
	delayedTaskMessage := "\t\t[Delay]"
	var dueDate time.Time
	var dueDateString string

	for _, list := range lists {

		var tasks []wunderlistutil.Task
		tasks = wunderlistutil.GetTasks(param, strconv.Itoa(list.ID), false)
		//sort by title=
		sort.Sort(ByTitle(tasks))
		taskString := "\t[" + list.Title + "]\n"
		applicableCount := 0
		//make message string
		for _, task := range tasks {
			dueDate, err = time.Parse(layout, task.DueDate)
			dueDateString = dueDate.Format(layout)
			if dueDate.Equal(today) {
				taskString += newTaskMessage + task.Title + "\n"
				applicableCount++
			} else if dueDate.Before(today) && (dueDateString != RESETED_DUE_DATE) {
				taskString += delayedTaskMessage + task.Title + "(" + dueDateString + ")\n"
				applicableCount++
			}
		}
		if applicableCount == 0 {
			continue
		} else {
			headMessage += taskString
		}
	}

	fmt.Println(headMessage)
	//	slackutil.sendmessage(headmessage + newtaskmessage + delayedtaskmessage)
}
