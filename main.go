package main

import (
	"fmt"
	"strings"
	"time"

	"gitlab.jasondale.me/jdale/govult/pkg/pidcheck"
	"gitlab.jasondale.me/jdale/govult/pkg/slack"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

// Config is the configuration struct
type config struct {
	ChannelToMessage string
	ChannelToMonitor string
	PidFilePath      string
	SlackBotToken    string
	SlackMessageText string
	SlackToken       string
	SlackUser        string
	SlackWebHook     string
	TriggerWords     []string
}

// new config instance
var (
	conf     *config
	SlackAPI slack.Service
)

func getConf() *config {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()

	if err != nil {
		handleError(err)
	}

	conf := &config{}
	err = viper.Unmarshal(conf)

	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}

//TODO: real error handling
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	conf = getConf()

	// * SlackAPI is still a Service from slack package; slack.Client satisfies the Service requirements,
	// * but SlackAPI will remain a slack.Service as it was declared up top.
	SlackAPI = slack.Client{
		ChannelToMessage: conf.ChannelToMessage,
		ChannelToMonitor: conf.ChannelToMonitor,
		SlackBotToken:    conf.SlackBotToken,
		SlackMessageText: conf.SlackMessageText,
		SlackToken:       conf.SlackToken,
		SlackWebHook:     conf.SlackWebHook,
	}
	pidPath := fmt.Sprintf("%s/goVult", conf.PidFilePath)
	pid := pidcheck.AlreadyRunning(pidPath)

	if !pid {
		// Infinite loop - get new messages every 5 seconds
		var LastMessageTs string
		LastMessageTs = "0"
		firstRun := true
		for true {
			messages, err := getSlackMessages(conf.ChannelToMonitor, LastMessageTs)
			if err != nil {
				fmt.Println("Error encountered: ", err)
			}
			for _, message := range messages.Messages {
				if message.Ts > LastMessageTs {
					SlackAPI.UpdateOldest(message.Ts)
					LastMessageTs = message.Ts
					spew.Dump(message)
				}
				if !firstRun {
					if len(message.Text) > 0 {
						if analyzeMessage(message.Text) {
							sendSlackMessage(message.Text)
						}

					}
				}
			}

			firstRun = false
			time.Sleep(5 * time.Second)
		}
	}
}

// Get Slack messages
func getSlackMessages(channel string, timestamp string) (slack.Response, error) {
	response, err := SlackAPI.GetMessages(channel, timestamp)

	return response, err
}

// Check a message for a match to any of the keywords
func analyzeMessage(message string) bool {
	match := false
	words := strings.Split(message, " ")
	for _, word := range words {
		for _, trigger := range conf.TriggerWords {
			if strings.Contains(word, trigger) {
				match = true
			}
		}
	}
	return match
}

// Send a slack message to a channel
func sendSlackMessage(message string) {
	SlackAPI.PostMessage(conf.SlackMessageText)
	SlackAPI.PostMessage("> " + message)
}
