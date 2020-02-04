package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// TODO: remove variables that are unused or need to be tracked elsewhere (oldest, channel stuff, etc)
// Client is the Slack Client
type Client struct {
	ChannelToMonitor string
	ChannelToMessage string
	client           *http.Client
	Oldest           string
	SlackBotToken    string
	SlackMessageText string
	SlackToken       string
	SlackUser        string
	SlackWebHook     string
	Token            string
}

// Error is the error
type Error struct {
	Ok           bool   `json:"ok"`
	ErrorText    string `json:"error"`
	CallResponse *http.Response
}

var baseURL = "https://slack.com/api"

func (e Error) Error() string {
	message := fmt.Sprintf("Slack API Error: %s \n Status Code: %d", e.ErrorText, e.CallResponse.StatusCode)
	return message
}

func (c *Client) call(method string, url string, payload interface{}, target interface{}) ([]byte, error) {
	url = baseURL + url
	if c.client == nil {
		c.client = &http.Client{}
	}
	// TODO: throw error if >= 400
	jsonData := []byte{}
	var err error

	if payload != nil {
		jsonData, err = json.Marshal(payload)
		if err != nil {
			return []byte{}, err
		}
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return []byte{}, err
	}

	request.Header.Add("Authorization", "Bearer "+c.SlackBotToken)
	request.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(request)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	//TODO: this can all be one error function, take responseBody and do all the error checks
	errorTarget := Error{}

	err = json.Unmarshal(responseBody, &errorTarget)
	if err != nil {
		return responseBody, err
	}

	if errorTarget.Ok != true {
		errorTarget.CallResponse = resp
		return responseBody, errorTarget
	}
	// TODO: ^ to here

	if resp.StatusCode >= 400 {
		err := fmt.Errorf("Slack HTTP Error: %d", resp.StatusCode)
		return responseBody, err
	}
	if target != nil {
		err = json.Unmarshal(responseBody, target)
		if err != nil {
			return responseBody, err
		}
	}

	return responseBody, nil
}
