package main

import (
	"fmt"
	"strconv"
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
	// sendSlackMessage("test")
	pidPath := fmt.Sprintf("%s/goVult", conf.PidFilePath)
	pid := pidcheck.AlreadyRunning(pidPath)
	var LastMessageTs int
	LastMessageTs = 0
	firstRun := true

	if !pid {
		// Infinite loop - get new messages every 5 seconds
		for {
			fmt.Println("Tick...")
			fmt.Println("LastMessageTs:", LastMessageTs)
			messages, err := getSlackMessages(conf.ChannelToMonitor, strconv.Itoa(LastMessageTs))
			// spew.Dump(messages)
			if err != nil {
				fmt.Println("Error encountered: ", err)
			}
			for _, message := range messages.Messages {
				currentTsSplit := strings.Split(message.Ts, ".")
				currentTs, err := strconv.Atoi(currentTsSplit[0])
				if err != nil {
					spew.Dump(err)
				}
				fmt.Printf("ts: %s,  %d, %d\n", message.Ts, LastMessageTs, currentTs)
				if currentTs > LastMessageTs {
					LastMessageTs = currentTs
					spew.Dump("Message: ", message)
					if !firstRun {
						if len(message.Text) > 0 {
							if analyzeMessage(message.Text) {
								sendSlackMessage(message.Text)
							}

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
				spew.Dump("We got a match", word)
			}
		}
	}
	return match
}

// Send a slack message to a channel
func sendSlackMessage(message string) {
	spew.Dump("Posting message: ", message)
	// err := SlackAPI.PostMessage(slack.Payload{
	// 	channel: conf.ChannelToMessage,
	// 	text:    message,
	// })
	// if err != nil {
	// 	spew.Dump(err)
	// 	// panic(err)
	// }
	// SlackAPI.PostMessage(conf.SlackMessageText)
	// SlackAPI.PostMessage("> " + message)
}
