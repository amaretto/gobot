package main

import (
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"../slackutil"

	"../wunderlistutil"
)

//for title Sort
type ByTitle []wunderlistutil.Task

func (a ByTitle) Len() int           { return len(a) }
func (a ByTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTitle) Less(i, j int) bool { return a[i].Title < a[j].Title }

// Messages
const beginMessage = "\n今日の実施タスクは...\n"
const noTaskMesssage = "一体一日何をしていたのですか？"
const fewTaskMessage = "タスク消化量が少なめです。もっと頑張ってください"
const manyTaskMessage = "素晴らしいタスクの消化量です！明日もこの調子でいきましょう"

func main() {

	param := wunderlistutil.Param{AccessToken: os.Getenv("WUND_ACTOKEN"),
		ClientID: os.Getenv("WUND_CLIENT")}

	var lists []wunderlistutil.List
	lists = wunderlistutil.GetLists(param)

	doneFlag := true

	//for date
	now := time.Now()
	today := now.Format("2006-01-02")
	message := beginMessage
	count := 0

	for _, list := range lists {
		var taskList []wunderlistutil.Task
		taskList = wunderlistutil.GetTasks(param, strconv.Itoa(list.ID), doneFlag)
		sort.Sort(ByTitle(taskList))
		applicableCount := 0
		taskString := ""
		for _, task := range taskList {
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
		message += noTaskMesssage
	} else if count < 5 {
		message += fewTaskMessage
	} else {
		message += manyTaskMessage
	}

	slackutil.SendMessage(message)
}
