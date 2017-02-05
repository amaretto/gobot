package slackutil

import (
	"os"

	"github.com/nlopes/slack"
)

func SendMessage(message string) {
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
