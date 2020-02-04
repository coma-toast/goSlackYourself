package slack

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
)

// Service is the Slack service
type Service interface {
	GetMessages(string, string) (Response, error)
	PostMessage(payload PostSlackMessage) error
}

// GetMessages gets Slack messages from a channel from a start time
func (c Client) GetMessages(channel string, timestamp string) (Response, error) {
	var response Response
	// messages, err := c.slackGet("conversations.history", channel, timestamp)
	_, err := c.call("GET", channel, timestamp, &response)
	spew.Dump("here", response)
	os.Exit(1)

	return Response{}, err
	// return messages, err
}

// GetUserInfo gets all user info for a userID
func (c Client) GetUserInfo(userid string) UserObject {
	var userData = UserObject{}

	return userData
}

// PostMessage posts a message to Slack
func (c Client) PostMessage(payload PostSlackMessage) error {
	url := "/chat.postMessage"

	response, err := c.call("POST", url, payload, nil)
	fmt.Println(string(response))

	return err
}
