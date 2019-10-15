package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	// "github.com/davecgh/go-spew/spew"
)

// Service is the Slack service
type Service interface {
	GetMessages() (Response, error)
	PostMessage(string) error
}

// Client is the Slack Client
type Client struct {
	ChannelToMonitor string
	ChannelToMessage string
	Oldest           string
	SlackBotToken    string
	SlackUser        string
	SlackWebHook     string
	Token            string
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
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

var baseURL = "https://slack.com/api"

// GetMessages gets Slack messages from a channel from a start time
func (c Client) GetMessages() (Response, error) {
	// messages, err := c.slackCall("GET", "conversations.history", c.ChannelToMonitor, c.Oldest, "")
	messages, err := c.slackGet("conversations.history", c.ChannelToMonitor, c.Oldest)

	return messages, err
}

// PostMessage posts a message to a channel
func (c Client) PostMessage(message string) error {
	// _, err := c.slackCall("POST", "/chat.postMessage", c.ChannelToMessage, "0", message)
	err := c.slackPost(message)

	return err
}

func (c Client) slackPost(message string) error {
	bearer := "Bearer " + c.SlackBotToken
	data := PostSlackMessage{
		Text:    message,
		Channel: c.ChannelToMessage,
	}
	// spew.Dump(data)
	// data.Set("text", message)
	// data.Set("channel", c.ChannelToMessage)
	jsonData, err := json.Marshal(data)

	client := &http.Client{}
	req, err := http.NewRequest("POST", c.SlackWebHook, bytes.NewBuffer(jsonData))
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")
	// req.Header.Add("Content-Length", strconv.Itoa(len(jsonData.Encode())))
	// ---
	// req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))
	// ---

	return err
}

func (c Client) slackGet(endpoint string, channel string, oldest string) (Response, error) {
	var messageData Response
	bearer := "Bearer " + c.SlackBotToken
	url := fmt.Sprintf("%s/%s", baseURL, endpoint)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	handleError(err)

	q := req.URL.Query()
	q.Add("channel", channel)
	q.Add("token", c.SlackBotToken)
	q.Add("oldest", oldest)
	// req, err := http.NewRequest("GET", url, strings.NewReader("token="+c.SlackBotToken))
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	handleError(err)

	err = json.Unmarshal(responseBody, &messageData)

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))
	// ---

	return messageData, err
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
	if method == "POST" {
		q.Add("as_user", "false")
		q.Add("username", c.SlackUser)
		q.Add("text", "test")
	}
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
