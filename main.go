package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"

	"gitlab.jasondale.me/jdale/govult/pkg/slack"

	"github.com/spf13/viper"
)

// Config is the configuration struct
type config struct {
	PidFilePath      string
	SlackToken       string
	ChannelToMonitor string
	ChannelToMessage string
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

// import config:
// token, channel to monitor, channel to message, trigger words
// pull messages since last pull from monitored channel
// see if any words are in the list from the message - some kind of string search
// if match, copy message and release the vultures
// rinse, repeat (while loop?)

func main() {
	conf = getConf()
	// * SlackAPI is still a Service from slack package; slack.Client satisfies the Service requirements,
	// * but SlackAPI will remain a slack.Service as it was declared up top.
	SlackAPI = slack.Client{
		Token:            conf.SlackToken,
		ChannelToMonitor: conf.ChannelToMonitor,
		ChannelToMessage: conf.ChannelToMessage,
		Oldest:           "0",
	}
	pidPath := fmt.Sprintf("%s/goVult", conf.PidFilePath)
	pid := alreadyRunning(pidPath)

	if !pid {
		// Infinite loop - get new messages every 5 seconds
		var LastMessageTs string
		LastMessageTs = "0"
		for true {
			messages, err := getNewSlackMessages()
			if err != nil {
				fmt.Println("Error encountered: ", err)
			}
			for _, message := range messages.Messages {
				if message.Ts > LastMessageTs {
					SlackAPI = slack.Client{
						Token:            conf.SlackToken,
						ChannelToMonitor: conf.ChannelToMonitor,
						ChannelToMessage: conf.ChannelToMessage,
						Oldest:           message.Ts,
					}
					LastMessageTs = message.Ts
				}
				// keywordMatchedMessages := analyzeMessage(message)
				// for _, matchedMessage := range keywordMatchedMessages {
				// 	fmt.Println("matchedMessage", matchedMessage)
				// sendSlackMessage(matchedMessage)
				// }
			}

			time.Sleep(5 * time.Second)
		}
	}
}

// Get Slack messages since last check
func getNewSlackMessages() (slack.Response, error) {
	response, err := SlackAPI.GetMessages()
	if len(response.Messages) > 0 {
		for _, message := range response.Messages {
			analyzeMessage(message.Text)
		}
	}
	return response, err
}

// Check a message for a match to any of the keywords
func analyzeMessage(message string) string {
	// fmt.Println("analyzeMessage " + message.Text)
	return message
}

// Send a slack message to a channel
func sendSlackMessage(message string) {
	// fmt.Println("sendSlackMessage " + message)
}
