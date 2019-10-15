package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Service is the Slack service
type Service interface {
	GetMessages() (Response, error)
	PostMessages([]string) error
}

// Client is the Slack Client
type Client struct {
	Token            string
	ChannelToMonitor string
	ChannelToMessage string
	Oldest           string
}

// Response is the all messages returned by the query
type Response struct {
	Ok       bool      `json:"ok"`
	Messages []Message `json:"messages"`
}

// Message is the individual message response
type Message struct {
	ClientMsgID string `json:"client_msg_id"`
	Type        string `json:"type"`
	Text        string `json:"text"`
	User        string `json:"user"`
	Ts          string `json:"ts"`
	Team        string `json:"team"`
}

// PostSlackMessage is the struct for posting messages to Slack
type PostSlackMessage struct {
	AsUser   string
	Channel  string
	Text     string
	Token    string
	Username string
}

var baseURL = "https://slack.com/api"

// GetMessages gets Slack messages from a channel from a start time
func (c Client) GetMessages() (Response, error) {
	messages, err := c.slackCall("GET", "conversations.history", c.ChannelToMonitor, c.Oldest, "")

	return messages, err
}

// PostMessages posts a slice of messages to a channel
func (c Client) PostMessages(messages []string) error {
	var err error
	for _, message := range messages {
		c.slackCall("POST", "/chat.postMessage", c.ChannelToMessage, "0", message)
	}

	return err
}

func (c Client) slackCall(method string, endpoint string, channel string, oldest string, data string) (Response, error) {
	var messageData Response
	url := fmt.Sprintf("%s/%s", baseURL, endpoint)
	client := &http.Client{}
	request, err := http.NewRequest(method, url, nil)
	handleError(err)

	q := request.URL.Query()
	q.Add("channel", channel)
	q.Add("token", c.Token)
	q.Add("oldest", oldest)
	q.
		request.URL.RawQuery = q.Encode()
	response, err := client.Do(request)
	handleError(err)

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	handleError(err)

	err = json.Unmarshal(responseBody, &messageData)

	if response.StatusCode >= 400 {
		log.Print(fmt.Errorf("Response Code Error: %d. %s", response.StatusCode, string(responseBody)))
	}

	return messageData, err
}

//TODO: real error handling
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
