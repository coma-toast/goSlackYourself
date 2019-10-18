package slack

import (
	"fmt"
)

// Service is the Slack service
type Service interface {
	GetMessages(string, string) (GetResponse, error)
	PostMessage(payload PostSlackMessage) (PostResponse, error)
}

// GetMessages gets Slack messages from a channel from a start time
func (c Client) GetMessages(channel string, timestamp string) (Response, error) {
	// messages, err := c.slackGet("conversations.history", channel, timestamp)
	messages, err := c.call("GET", channel, timestamp)

	return messages, err
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
