package main

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

func basicChannelExample() {
	// create two channels
	ping := make(chan string)
	pong := make(chan string)
	go shout(ping, pong)

	fmt.Println("Type something and press enter (enter Q to quit)")

	for {
		color.Cyan("--> ")
		// get user input
		var input string
		_, _ = fmt.Scanln(&input)

		if input == strings.ToLower("q") {
			break
		}
		ping <- input

		// wait for response
		response := <-pong

		color.Green("Response: " + response)
	}

	color.Red("Exiting...")
	close(ping)
	close(pong)

}

func shout(ping <-chan string, pong chan<- string) {
	for {
		s := <-ping
		pong <- fmt.Sprintf("%s !!! ", strings.ToUpper(s))
	}
}
