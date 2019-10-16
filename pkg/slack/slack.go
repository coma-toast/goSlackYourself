package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Service is the Slack service
type Service interface {
	GetMessages(string, string) (Response, error)
	PostMessage(string) error
}

// Client is the Slack Client
type Client struct {
	ChannelToMonitor string
	ChannelToMessage string
	Oldest           string
	SlackBotToken    string
	SlackMessageText string
	SlackToken       string
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

// UserObject is the struct for user data
type UserObject struct {
	ID       string `json:"id"`
	TeamID   string `json:"team_id"`
	Name     string `json:"name"`
	Deleted  bool   `json:"deleted"`
	Color    string `json:"color"`
	RealName string `json:"real_name"`
	Tz       string `json:"tz"`
	TzLabel  string `json:"tz_label"`
	TzOffset int    `json:"tz_offset"`
	Profile  struct {
		AvatarHash            string `json:"avatar_hash"`
		StatusText            string `json:"status_text"`
		StatusEmoji           string `json:"status_emoji"`
		RealName              string `json:"real_name"`
		DisplayName           string `json:"display_name"`
		RealNameNormalized    string `json:"real_name_normalized"`
		DisplayNameNormalized string `json:"display_name_normalized"`
		Email                 string `json:"email"`
		ImageOriginal         string `json:"image_original"`
		Image24               string `json:"image_24"`
		Image32               string `json:"image_32"`
		Image48               string `json:"image_48"`
		Image72               string `json:"image_72"`
		Image192              string `json:"image_192"`
		Image512              string `json:"image_512"`
		Team                  string `json:"team"`
	} `json:"profile"`
	IsAdmin           bool `json:"is_admin"`
	IsOwner           bool `json:"is_owner"`
	IsPrimaryOwner    bool `json:"is_primary_owner"`
	IsRestricted      bool `json:"is_restricted"`
	IsUltraRestricted bool `json:"is_ultra_restricted"`
	IsBot             bool `json:"is_bot"`
	Updated           int  `json:"updated"`
	IsAppUser         bool `json:"is_app_user"`
	Has2Fa            bool `json:"has_2fa"`
}

var baseURL = "https://slack.com/api"

// GetMessages gets Slack messages from a channel from a start time
func (c Client) GetMessages(channel string, timestamp string) (Response, error) {
	messages, err := c.slackGet("conversations.history", channel, timestamp)

	return messages, err
}

// GetUserInfo gets all user info for a userID
func (c Client) GetUserInfo(userid string)

// PostMessage posts a message to a channel
func (c Client) PostMessage(message string) error {
	err := c.slackPost(message)
	return err
}

func (c Client) slackPost(message string) error {
	bearer := "Bearer " + c.SlackBotToken
	data := PostSlackMessage{
		Text:    message,
		Channel: c.ChannelToMessage,
	}
	jsonData, err := json.Marshal(data)

	client := &http.Client{}
	req, err := http.NewRequest("POST", c.SlackWebHook, bytes.NewBuffer(jsonData))
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return err
}

func (c Client) slackGet(endpoint string, channel string, oldest string) (Response, error) {
	var messageData Response
	url := fmt.Sprintf("%s/%s", baseURL, endpoint)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	handleError(err)

	q := req.URL.Query()
	q.Add("channel", channel)
	q.Add("token", c.SlackToken)
	q.Add("oldest", oldest)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	handleError(err)

	err = json.Unmarshal(responseBody, &messageData)

	return messageData, err
}

// TODO: replace this function with parts of the slackGet and slackPost so we are DRY
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
