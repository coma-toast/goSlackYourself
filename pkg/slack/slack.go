package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Service is the Slack service
type Service interface {
	GetMessages(string, string) (GetResponse, error)
	PostMessage(payload PostSlackMessage) (PostResponse, error)
}

// GetMessages gets Slack messages from a channel from a start time
func (c Client) GetMessages(channel string, timestamp string) (Response, error) {
	messages, err := c.slackGet("conversations.history", channel, timestamp)

	return messages, err
}

// GetUserInfo gets all user info for a userID
func (c Client) GetUserInfo(userid string) UserObject {
	return UserObject{}
}

// PostMessage posts a message to Slack
func (c Client) PostMessage(payload PostSlackMessage) error {
	url := "/chat.postMessage"

	response, err := c.call("POST", url, payload, nil)
	fmt.Println(string(response))

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
	q.Add("oldest", oldest)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	handleError(err)

	err = json.Unmarshal(responseBody, &messageData)

	return messageData, err
}
