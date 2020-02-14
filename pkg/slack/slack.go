package slack

import (
	"fmt"
)

// Service is the Slack service
type Service interface {
	GetMessages(string, string) (Response, error)
	PostMessage(payload Payload) error
}

// GetMessages gets Slack messages from a channel from a start time
func (c Client) GetMessages(channel string, timestamp string) (Response, error) {
	var response Response
	payload := Payload{
		channel: channel,
		oldest:  timestamp,
		token:   c.SlackBotToken,
	}
	// messages, err := c.slackGet("conversations.history", channel, timestamp)
	_, err := c.call("GET", "channels.history", payload, &response)
	// spew.Dump("response: ", response)
	// os.Exit(1)

	return response, err
	// return messages, err
}

// GetUserInfo gets all user info for a userID
func (c Client) GetUserInfo(userid string) UserObject {
	var userData = UserObject{}

	return userData
}

// PostMessage posts a message to Slack
func (c Client) PostMessage(payload Payload) error {
	url := "/chat.postMessage"

	response, err := c.call("POST", url, payload, nil)
	fmt.Println(string(response))

	return err
}
