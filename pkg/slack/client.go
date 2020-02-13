package slack

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

var baseURL = "https://slack.com/api/"

func (e Error) Error() string {
	message := fmt.Sprintf("Slack API Error: %s \n Status Code: %d", e.ErrorText, e.CallResponse.StatusCode)
	return message
}

func (c *Client) call(method string, destination string, payload Payload, target interface{}) ([]byte, error) {
	destination = baseURL + destination
	jsonData := []byte{}
	_ = jsonData
	var err error

	if c.client == nil {
		c.client = &http.Client{}
	}

	values := url.Values{
		"token":   {payload.token},
		"channel": {payload.channel},
	}

	req, err := http.NewRequest(method, destination, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = values.Encode()

	// spew.Dump("Request: ", req)

	resp, err := c.client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	// spew.Dump("ReponseBody: ", responseBody)
	//TODO: this can all be one error function, take responseBody and do all the error checks
	errorTarget := Error{}

	// err = json.Unmarshal(responseBody, &target)
	// if err != nil {
	// 	spew.Dump("failed to unmarshal json")
	// 	return responseBody, err
	// }

	// spew.Dump("ErrorTarget: ", errorTarget)

	if errorTarget.Ok != true {
		errorTarget.CallResponse = resp
		return responseBody, errorTarget
	}
	// TODO: ^ to here

	if resp.StatusCode >= 400 {
		err := fmt.Errorf("Slack HTTP Error: %d", resp.StatusCode)
		return responseBody, err
	}
	// if target != nil {
	// 	err = json.Unmarshal(responseBody, target)
	// 	if err != nil {
	// 		return responseBody, err
	// 	}
	// }

	return responseBody, nil
}
