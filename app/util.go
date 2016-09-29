package app

import (
	"math/rand"
	"regexp"
	"strings"
)

var insults = [...]string{"Stupid.", "You idiot.", "You fool."}

func init() {
	rand.Seed(7483658374658473)
}

func matches(keyword string, command string) bool {
	return len(keyword) > 0 && len(keyword) <= len(command) && command[:len(keyword)] == keyword
}

func randomInsult() string {
	return insults[rand.Intn(len(insults))]
}

func ReadNextCommand(input string) (keyword string, newInput string) {
	keyword = strings.ToLower(input)

	re := regexp.MustCompile("\\s+")
	loc := re.FindStringIndex(input)
	if loc != nil {
		keyword = strings.ToLower(input[:loc[0]])
		newInput = input[loc[1]:]
	}

	return
}
