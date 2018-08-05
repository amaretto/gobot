package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nlopes/slack"
)

func HandlerRequest(ctx context.Context, params interface{}) (interface{}, error) {
	SendMessage("hello")
	return params, nil
}

func SendMessage(message string) {
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	username := os.Getenv("SLACK_USERNAME")

	message = "<" + username + "> " + message
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Text: message,
	}
	params.Attachments = []slack.Attachment{attachment}
	params.AsUser = true

	api.PostMessage("#todo", "", params)
}

func main() {
	lambda.Start(HandlerRequest)
}
