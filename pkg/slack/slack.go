package slack

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

// Service is the Slack service
type Service interface {
	GetMessages(channel string, timestamp string) (Response, error)
	GetUserInfo(userid string) Response
	PostMessage(channel string, text string) error
}

// GetMessages gets Slack messages from a channel from a start time
func (c Client) GetMessages(channel string, timestamp string) (Response, error) {
	var response Response
	payload := Payload{
		channel: channel,
		oldest:  timestamp,
		token:   c.SlackBotToken,
	}

	err := c.call("GET", "channels.history", payload, &response)

	return response, err
}

// GetUserInfo gets all user info for a userID
func (c Client) GetUserInfo(userid string) Response {
	var userData = Response{}
	url := "users.info"
	payload := Payload{
		token: c.SlackBotToken,
		user:  userid,
	}

	err := c.call("GET", url, payload, &userData)
	spew.Dump(userData)

	if err != nil {
		fmt.Printf("error getting user info: %e", err)
	}

	return userData
}

// PostMessage posts a message to Slack
func (c Client) PostMessage(channel string, text string) error {
	url := "chat.postMessage"
	payload := Payload{
		channel:    channel,
		token:      c.SlackBotToken,
		icon_emoji: ":vulture:",
		text:       text,
	}
	var response Response

	err := c.call("POST", url, payload, &response)

	return err
}
