package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

// Service is the Slack service
type Service interface {
	GetMessages() ([]string, error)
	PostMessages([]string) error
}

// Client is the Slack Client
type Client struct {
	Token            string
	ChannelToMonitor string
	ChannelToMessage string
	Oldest           int64
}

var baseURL = "https://slack.com/api/"

// GetMessages gets Slack messages from a channel from a start time
func (c Client) GetMessages() ([]string, error) {
	messages, err := c.slackCall("GET", "conversations.history", c.ChannelToMonitor, c.Oldest, "")
	spew.Dump("GetMessages ", messages)

	return messages, err
}

// PostMessages posts a slice of messages to a channel
func (c Client) PostMessages(messages []string) error {
	var err error
	for _, message := range messages {
		c.slackCall("POST", "/chat.postMessage", c.ChannelToMessage, 0, message)
	}

	return err
}

func (c Client) slackCall(method string, endpoint string, channel string, startTime int64, data string) ([]string, error) {
	var messageData []string
	url := fmt.Sprintf("%s/%s", baseURL, endpoint)
	client := &http.Client{}
	request, err := http.NewRequest(method, url, nil)
	handleError(err)
	// add authorization header to the request

	request.Header.Add("token", c.Token)
	response, err := client.Do(request)
	handleError(err)
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	handleError(err)
	err = json.Unmarshal(responseBody, &messageData)
	// if response.StatusCode >= 400 {
	log.Print(fmt.Errorf("Response Code Error: %d. %s", response.StatusCode, string(responseBody)))

	// }
	return messageData, err
}

//TODO: real error handling
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
