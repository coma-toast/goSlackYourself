package main

func main() {
	words := getListOfTriggerWords()
	messages := getSlackMessages()
	food := analyzeMessages(words, messages)
	releaseTheVultures(food)
}
