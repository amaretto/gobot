package main

import (
	"os"

	"github.com/nlopes/slack"
)

func SendMessage(message string) {
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	username := os.Getenv("SLACK_USERNAME")

	pretext := username + " " + message

	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Pretext: pretext,
	}
	params.Attachments = []slack.Attachment{attachment}
	params.AsUser = true
	params.LinkNames = 1

	api.PostMessage("#todo", "", params)
}

func main() {
	SendMessage("hoge #todo")
}
