package slack

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Slack service {
	GetMessages string
	PostMessages string
}

func slackCall(method string, endpoint string, channel string, data string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", conf.BaseURL, endpoint)
	client := &http.Client{}
	request, err := http.NewRequest(method, url, nil)
	handleError(err)
	// add authorization header to the request
	request.Header.Add("Authorization", conf.SlackToken)
	response, err := client.Do(request)
	handleError(err)
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	handleError(err)
	// if response.StatusCode >= 400 {
	log.Print(fmt.Errorf("Response Code Error: %d. %s", response.StatusCode, string(data)))

	// }
	return data, err
}
