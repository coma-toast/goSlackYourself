package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"

	"gitlab.jasondale.me/jdale/govult/slack"

	"github.com/spf13/viper"
)

// config is the configuration struct
type config struct {
	PidFilePath      string
	SlackToken       string
	BaseURL          string
	ChannelToMonitor string
	ChannelToMessage string
	TriggerWords     []string
}

// new config instance
var (
	conf *config
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
	pidPath := fmt.Sprintf("%s/goVult", conf.PidFilePath)
	pid := alreadyRunning(pidPath)

	if !pid {
		// Infinite loop - get new messages every 5 seconds
		for true {
			messages := getNewSlackMessages()
			for message := range messages {
				keywordMatchedMessages := analyzeMessage(message)
				for matchedMessage := range keywordMatchedMessages {
					sendSlackMessage(matchedMessage)
				}
			}
			time.Sleep(5 * time.Second)
		}
	}
}

// Get Slack messages since last check
func getNewSlackMessages() string {
	messages := slack.GetMessages("GET", "conversations.history", conf.ChannelToMonitor, nil)
	return messages
}

// Check a message for a match to any of the keywords
func analyzeMessage(message string) string {
	fmt.Println("analyzeMessage" + message)
	return message
}

// Send a slack message to a channel
func sendSlackMessage(message string, channel string) {
	fmt.Println("sendSlackMessage" + message)
}
