package main

import (
	"fmt"
	"github.com/fatih/color"
	"time"
)

type Philosopher struct {
	name string
}

var philosophers = []Philosopher{
	{name: "Plato"},
	{name: "Socrates"},
	{name: "Aristotle"},
	{name: "Pascal"},
	{name: "Locke"},
}

var hunger = 3 // The number of times each philosopher eats
var eatTime = 1 * time.Second
var thinkTime = 1 * time.Second

func dine() {
	// Create channels to represent forks.
	numPhilosophers := len(philosophers)
	forks := make([]chan struct{}, numPhilosophers)

	// Initialize each fork as a channel with a buffer size of 1.
	for i := 0; i < numPhilosophers; i++ {
		forks[i] = make(chan struct{}, 1)
		// Add an initial value to each fork's channel to represent that it's available.
		forks[i] <- struct{}{}
	}

	// Channel to signal when a philosopher finishes dining.
	done := make(chan string, numPhilosophers)

	// Start the dining process for each philosopher.
	for i := 0; i < numPhilosophers; i++ {
		if i%2 == 0 {
			// Even-indexed philosophers pick up the left fork first
			go philosopherRoutine(philosophers[i], forks[i], forks[(i+1)%numPhilosophers], done, true)
		} else {
			// Odd-indexed philosophers pick up the right fork first
			go philosopherRoutine(philosophers[i], forks[(i+1)%numPhilosophers], forks[i], done, false)
		}
	}

	// Wait for all philosophers to finish.
	for i := 0; i < numPhilosophers; i++ {
		color.Green(<-done)
	}

	color.Cyan("--------------------------")
	color.Cyan("Philosophers are happy and full!")
}

func philosopherRoutine(philosopher Philosopher, firstFork, secondFork chan struct{}, done chan string, isLeftFirst bool) {
	for i := hunger; i > 0; i-- {
		if isLeftFirst {
			color.Yellow("> %s tries to take the left fork", philosopher.name)
		} else {
			color.Yellow("> %s tries to take the right fork", philosopher.name)
		}

		// Pick up the first fork.
		<-firstFork
		if isLeftFirst {
			color.Yellow("> %s takes the left fork", philosopher.name)
			color.Yellow("> %s tries to take the right fork", philosopher.name)
		} else {
			color.Yellow("> %s takes the right fork", philosopher.name)
			color.Yellow("> %s tries to take the left fork", philosopher.name)
		}

		// Pick up the second fork.
		<-secondFork

		color.Cyan("%s has both forks and is eating", philosopher.name)
		time.Sleep(eatTime)

		// Put down the forks by sending back to the channels.
		firstFork <- struct{}{}
		secondFork <- struct{}{}
		color.Green("%s has put down both forks", philosopher.name)

		// Think
		color.Blue("%s is thinking...", philosopher.name)
		time.Sleep(thinkTime)
	}

	// Philosopher is finished.
	done <- fmt.Sprintf("%s is finished dining and is sleeping...", philosopher.name)
}

func main() {
	color.Cyan("Dining Philosopher Problem")
	color.Cyan("--------------------------")
	color.HiYellow("The table is empty and the philosophers are hungry")

	dine()
}
