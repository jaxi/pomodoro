package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	WORKING_TIME = 25
	SHORT_BREAK  = 5
	LONG_BREAK   = 15
	PERIOD_COUNT = 4
)

type SpotifyCommand struct {
	Statement string
	Value     string
}

var cmdMap = map[string]string{
	"play":       "play",
	"pause":      "pause",
	"next_track": "next track",
	"play_track": "play track",
}

func (c *SpotifyCommand) Run() (string, error) {
	cmdName := "/usr/bin/osascript"
	spfx := `tell application "Spotify" to`

	cmdSlice := []string{spfx, cmdMap[c.Statement], c.Value}
	osaScript := strings.Join(cmdSlice, " ")

	output, err := exec.Command(cmdName, "-e", osaScript).Output()
	if err != nil {
		log.Fatal(err)
		log.Fatal(output)

		return "", err
	}

	result := bytes.TrimSpace(output)
	return string(result), nil
}

var (
	BreakMusic  SpotifyCommand = SpotifyCommand{Statement: "play_track", Value: `"spotify:track:3UQM3V4mjS1DuAqucivt1Q"`}
	StopPlaying SpotifyCommand = SpotifyCommand{Statement: "pause"}
)

var timeInterval time.Duration

func init() {
	debug := os.Getenv("DEBUG")

	if debug != "" {
		timeInterval = time.Second
	} else {
		timeInterval = time.Minute
	}
}

func _break(length int, breakEndMessage string) {
	ch := make(chan int)

	BreakMusic.Run()

	select {
	case <-ch:
	case <-time.After(time.Duration(length) * timeInterval):
		StopPlaying.Run()
		log.Println(breakEndMessage)
	}
}

func shortBreak() {
	log.Print("Short break started")
	_break(SHORT_BREAK, "Short break time out")
}

func longBreak() {
	log.Println("Long break started")
	_break(LONG_BREAK, "Long break time out")
}

func working() {
	log.Println("Started Working...")
	time.Sleep(WORKING_TIME * timeInterval)
}
func run() {
	log.Println("Pomodoro started")

	for {
		for i := 1; i <= PERIOD_COUNT; i++ {
			working()
			if i == PERIOD_COUNT {
				longBreak()
			} else {
				shortBreak()
			}
		}
	}
}

func main() {
	run()
}
