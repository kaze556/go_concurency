package main

import (
	"fmt"
	"github.com/fatih/color"
	"sync"
	"time"
)

// Philosopher info about philosopher
type Philosopher struct {
	name                string
	rightFork, leftFork int
}

// list of all Philosopher
var philosophers = []Philosopher{
	{name: "Plato", rightFork: 4, leftFork: 0},
	{name: "Socrates", rightFork: 0, leftFork: 1},
	{name: "Aristotle", rightFork: 1, leftFork: 2},
	{name: "Pascal", rightFork: 2, leftFork: 3},
	{name: "Locke", rightFork: 3, leftFork: 4},
}

var hunger = 3 // times a pearson eat
var eatTime = 1 * time.Second
var thinkTime = 1 * time.Second

var orderLock = &sync.Mutex{}
var orderFinished []string

func dine() {
	// eating wait group is needed because we need to wait for all philosophers to finish eating before the program terminates
	eatingWg := &sync.WaitGroup{}
	eatingWg.Add(len(philosophers))

	// seated wait group is needed because we need to wait for all philosophers to be seated before eating
	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	// forks is map of mutex because each fork (shared resource) needs concurrent access control
	var forks = make(map[int]*sync.Mutex)

	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	// start the meal
	for i := 0; i < len(philosophers); i++ {
		go tryEating(philosophers[i], eatingWg, forks, seated)
	}

	eatingWg.Wait()
}

func tryEating(
	philosopher Philosopher,
	eatingWg *sync.WaitGroup,
	forks map[int]*sync.Mutex,
	seated *sync.WaitGroup) {
	defer eatingWg.Done()
	// seat the philosopher
	fmt.Println(philosopher.name, "is seated")
	seated.Done()

	// we do this because seated wg ensures that all philosophers are seated before they start eating,
	// preventing any philosopher from prematurely attempting to eat before others are ready
	seated.Wait()
	// eat 3 times or (hunger times)
	for i := hunger; i > 0; i-- {

		if philosopher.leftFork > philosopher.rightFork {
			forks[philosopher.rightFork].Lock()
			color.Yellow("> %s takes a right fork", philosopher.name)
			forks[philosopher.leftFork].Lock()
			color.Yellow("> %s takes a left fork", philosopher.name)
		} else {
			forks[philosopher.leftFork].Lock()
			color.Yellow("> %s takes a left fork", philosopher.name)
			forks[philosopher.rightFork].Lock()
			color.Yellow("> %s takes a right fork", philosopher.name)
		}

		color.Cyan("%s has both forks and is eating", philosopher.name)
		time.Sleep(eatTime)

		color.Blue("%s is done eating and is thinking...", philosopher.name)
		time.Sleep(thinkTime)

		forks[philosopher.rightFork].Unlock()
		forks[philosopher.leftFork].Unlock()

		color.Green("%s has put down both forks", philosopher.name)
	}

	sprintf := fmt.Sprintf("				%s is finished and is sleeping...", philosopher.name)
	orderLock.Lock()
	orderFinished = append(orderFinished, sprintf)
	orderLock.Unlock()
}

func diningExample() {
	color.Cyan("Dining Philosopher Problem")
	color.Cyan("--------------------------")
	color.HiYellow("The table is empty and the philosophers are hungry")

	// start the meal
	dine()

	for _, message := range orderFinished {
		color.Red(message)
	}
	// finished message
	color.Green("--------------------------")
	color.Green("Philosophers are happy and Full")
}
