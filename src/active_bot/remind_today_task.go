package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"

	"../slackutil"

	"encoding/json"

	"../wunderlistutil"
)

//for title Sort
type ByTitle []wunderlistutil.Task

func (a ByTitle) Len() int           { return len(a) }
func (a ByTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTitle) Less(i, j int) bool { return a[i].Title < a[j].Title }

//for avoiding invalid due date
const RESETED_DUE_DATE = "0001-01-01"

// for checking counts difference from yesterday
type YstCounts struct {
	Today  int
	Delay  int
	Future int
}

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
	headMessage := "本日以降の実施予定タスクは以下です\n"
	newTaskMessage := "\t\t[New]"
	delayedTaskMessage := "\t\t[Delay]"
	futureTaskMessage := "\t\t[Future]"
	var dueDate time.Time
	var dueDateString string

	totalTodayCount := 0
	totalDelayCount := 0
	totalFutureCount := 0

	for _, list := range lists {

		var tasks []wunderlistutil.Task
		tasks = wunderlistutil.GetTasks(param, strconv.Itoa(list.ID), false)
		//sort by title=
		sort.Sort(ByTitle(tasks))
		taskString := "\t[" + list.Title + "]===================\n"
		applicableCount := 0
		//make message string
		for _, task := range tasks {
			dueDate, err = time.Parse(layout, task.DueDate)
			dueDateString = dueDate.Format(layout)
			if dueDate.Equal(today) {
				taskString += newTaskMessage + task.Title + "\n"
				totalTodayCount++
				applicableCount++
			} else if dueDate.Before(today) && (dueDateString != RESETED_DUE_DATE) {
				taskString += delayedTaskMessage + task.Title + "(" + dueDateString + ")\n"
				totalDelayCount++
				applicableCount++
			} else if dueDateString != RESETED_DUE_DATE {
				taskString += futureTaskMessage + task.Title + "(" + dueDateString + ")\n"
				totalFutureCount++
				applicableCount++
			}
		}
		if applicableCount == 0 {
			continue
		} else {
			headMessage += taskString
		}
	}

	// load yesterday count
	var yst YstCounts
	raw, err := ioutil.ReadFile("./hoge")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	json.Unmarshal(raw, &yst)

	// formatting message
	headMessage += "=================="
	headMessage += "\n今日が期限のタスク:\t" + strconv.Itoa(totalTodayCount)
	if yst.Today > totalTodayCount {
		headMessage += ":arrow_lower_right:"
	} else if yst.Today < totalTodayCount {
		headMessage += ":arrow_upper_right:"
	} else {
		headMessage += ":arrow_right:"
	}

	headMessage += "\n期限切れのタスク:\t" + strconv.Itoa(totalDelayCount)
	if yst.Delay > totalDelayCount {
		headMessage += ":arrow_lower_right:"
	} else if yst.Delay < totalDelayCount {
		headMessage += ":arrow_upper_right:"
	} else {
		headMessage += ":arrow_right:"
	}

	headMessage += "\n明日以降期限のタスク:\t" + strconv.Itoa(totalFutureCount)
	if yst.Future > totalFutureCount {
		headMessage += ":arrow_lower_right:"
	} else if yst.Future < totalFutureCount {
		headMessage += ":arrow_upper_right:"
	} else {
		headMessage += ":arrow_right:"
	}

	// register today's counts for file
	yst.Today = totalTodayCount
	yst.Delay = totalDelayCount
	yst.Future = totalFutureCount
	jsonBytes, err := json.Marshal(yst)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return
	}
	file, err := os.Create("./hoge")
	if err != nil {
		fmt.Println("File Create error:", err)
		return
	}
	file.Write(jsonBytes)

	slackutil.SendMessage(headMessage)
	//fmt.Println(headMessage)
}
