package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gitlab.jasondale.me/jdale/govult/pkg/slack"

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

// Write a pid file, but first make sure it doesn't exist with a running pid.
func alreadyRunning(pidFile string) bool {
	// Read in the pid file as a slice of bytes.
	if piddata, err := ioutil.ReadFile(pidFile); err == nil {
		// Convert the file contents to an integer.
		if pid, err := strconv.Atoi(string(piddata)); err == nil {
			// Look for the pid in the process list.
			if process, err := os.FindProcess(pid); err == nil {
				// Send the process a signal zero kill.
				if err := process.Signal(syscall.Signal(0)); err == nil {
					fmt.Println("PID already running!")
					// We only get an error if the pid isn't running, or it's not ours.
					err := fmt.Errorf("pid already running: %d", pid)
					log.Print(err)
					return true
				}
				log.Print(err)

			} else {
				log.Print(err)
			}
		} else {
			log.Print(err)
		}
	} else {
		log.Print(err)
	}
	// If we get here, then the pidfile didn't exist,
	// or the pid in it doesn't belong to the user running this app.
	ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0664)
	return false
}

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
		Oldest:           "0",
		SlackToken:       conf.SlackToken,
		SlackMessageText: conf.SlackMessageText,
		SlackWebHook:     conf.SlackWebHook,
	}
	pidPath := fmt.Sprintf("%s/goVult", conf.PidFilePath)
	pid := alreadyRunning(pidPath)

	if !pid {
		// Infinite loop - get new messages every 5 seconds
		var LastMessageTs string
		LastMessageTs = "0"
		firstRun := true
		for true {
			messages, err := getSlackMessages()
			if err != nil {
				fmt.Println("Error encountered: ", err)
			}
			for _, message := range messages.Messages {
				if message.Ts > LastMessageTs {
					SlackAPI = slack.Client{
						ChannelToMonitor: conf.ChannelToMonitor,
						ChannelToMessage: conf.ChannelToMessage,
						Oldest:           message.Ts,
						SlackBotToken:    conf.SlackBotToken,
						SlackMessageText: conf.SlackMessageText,
						SlackToken:       conf.SlackToken,
						SlackUser:        conf.SlackUser,
						SlackWebHook:     conf.SlackWebHook,
						Token:            conf.SlackToken,
					}
					LastMessageTs = message.Ts
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
func getSlackMessages() (slack.Response, error) {
	response, err := SlackAPI.GetMessages()

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
