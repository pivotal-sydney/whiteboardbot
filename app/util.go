package app
import (
	"strings"
	"regexp"
	"math/rand"
	. "github.com/pivotal-sydney/whiteboardbot/model"
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

func readNextCommand(input string) (keyword string, newInput string) {
	re := regexp.MustCompile("\\s+")
	loc := re.FindStringIndex(input)
	if loc != nil {
		keyword = strings.ToLower(input[:loc[0]])
		newInput = input[loc[1]:]
	} else {
		keyword = strings.ToLower(input)
		newInput = ""
	}
	return
}

func missingEntry(entryType EntryType) bool {
	return entryType == nil
}