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
		SlackUser:        conf.SlackUser,
		SlackWebHook:     conf.SlackWebHook,
	}
	pidPath := fmt.Sprintf("%s/goVult", conf.PidFilePath)
	pid := pidcheck.AlreadyRunning(pidPath)
	var LastMessageTs int
	LastMessageTs = 0
	firstRun := true

	if !pid {
		// Infinite loop - get new messages every 5 seconds
		for {
			// fmt.Println("Tick...") // * dev code
			messages, err := getSlackMessages(conf.ChannelToMonitor, strconv.Itoa(LastMessageTs))
			if err != nil {
				fmt.Println("Error encountered: ", err)
			}
			for _, message := range messages.Messages {
				currentTsSplit := strings.Split(message.Ts, ".")
				currentTs, err := strconv.Atoi(currentTsSplit[0])
				if err != nil {
					spew.Dump(err)
				}
				if currentTs > LastMessageTs {
					LastMessageTs = currentTs
					if !firstRun {
						if len(message.Text) > 0 {
							if analyzeMessage(message.Text) {
								sendSlackMessage(message)
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
	words := strings.Split(message, " ")
	for _, word := range words {
		word = strings.ToLower(word)
		for _, trigger := range conf.TriggerWords {
			if strings.Contains(word, trigger) {
				return true
			}
		}
	}
	return false
}

// Send a slack message to a channel
func sendSlackMessage(message slack.Message) {
	userData := SlackAPI.GetUserInfo(message.User)
	timestampSplit := strings.Split(message.Ts, ".")
	timestampInt, err := strconv.ParseInt(timestampSplit[0], 10, 64)
	timestamp := time.Unix(timestampInt, 0)
	if err != nil {
		fmt.Println(err)
	}

	err = SlackAPI.PostMessage(conf.ChannelToMessage, conf.SlackMessageText)
	if err != nil {
		spew.Dump(err)
	}

	err = SlackAPI.PostMessage(conf.ChannelToMessage, "> <@"+userData.User.ID+"> - "+timestamp.Format("03:04:05 PM")+": \n> "+message.Text)
	if err != nil {
		spew.Dump(err)
	}
}
