package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"

	"github.com/spf13/viper"
)

// config is the configuration struct
type config struct {
	BaseURL     string
	Words       []string
	PidFilePath string
	StartTime   string
	SlackToken  string
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

// ! example start
	
		// Start at the root folder of your choosing (i.e. Completed),
		// recursively searching down, populating the files list
	// 	files, err := getFilesFromFolder(conf.CompletedFolder)
	// 	if err != nil {
	// 		fmt.Println(fmt.Errorf("Error: %s", err))
	// 		return
	// 	}
	// 	downloadFiles(files)
	// 	deleteDownloaded(DeleteQueue)
	// }
// }

// func apiCall(method string, id int, callType string) ([]byte, error) {
// 	url := fmt.Sprintf("%s/%s", conf.BaseURL, callType)
// 	if id != 0 {
// 		url = fmt.Sprintf("%s/%d", url, id)
// 	}
// 	client := &http.Client{}
// 	request, err := http.NewRequest(method, url, nil)
// 	handleError(err)
// 	request.SetBasicAuth(conf.Username, conf.Passwd)
// 	response, err := client.Do(request)
// 	handleError(err)
// 	defer response.Body.Close()
// 	data, err := ioutil.ReadAll(response.Body)
// 	handleError(err)
// 	// if response.StatusCode >= 400 {
// 	log.Print(fmt.Errorf("Response Code Error: %d. %s", response.StatusCode, string(data)))

// 	// }
// 	return data, err
// }

//TODO: once errors are returned above, this is not needed
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// ! example end


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
	messages := []string

	// Infinite loop - get new messages every 5 seconds
	for {
		getNewSlackMessages()
		time.Sleep(5 * time.Second)
	}

	if !pid {}

	words := getListOfTriggerWords()
	food := analyzeMessages(words, messages)
	releaseTheVultures(food)
}
