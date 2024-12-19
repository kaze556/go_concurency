package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"strconv"
	"time"
)

var numberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

// PizzaOrder represents a single orderFinished for a pizza, including details like the orderFinished number,
// a message explaining its status, and whether the orderFinished was successfully completed.
type PizzaOrder struct {
	pizzaNumber int    // The number of the pizza orderFinished
	message     string // A message describing the status of the pizza orderFinished
	success     bool   // Whether the pizza orderFinished was successfully completed
}

// Producer is responsible for creating pizza orders and manages the flow of data via channels.
// It includes a data channel for sending pizza orders and a quit channel to signal stopping production.
type Producer struct {
	data chan PizzaOrder // Channel to send pizza orders to the consumer
	quit chan chan error // Channel to signal shutting down production
}

// Close gracefully closes the producer by signaling via the quit channel and waiting for confirmation.
func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

// pizzeria is the producer function that generates pizza orders continuously until stopped.
// It ensures each orderFinished is processed and sends either the completed or failed pizza to the consumer.
func pizzeria(pizzaMaker *Producer) {
	i := 0 // Keeps track of the current pizza number
	for {
		pizza := makePizza(i) // Try to make a pizza
		if pizza != nil {
			i = pizza.pizzaNumber // Update the pizza number
			select {
			case pizzaMaker.data <- *pizza: // Send the pizza orderFinished to the data channel
			case quitChan := <-pizzaMaker.quit: // Handle a quit signal
				close(pizzaMaker.data) // Close the data channel
				close(quitChan)        // Acknowledge the quit signal
				return
			}
		}
	}
}

// makePizza tries to create a single pizza orderFinished and returns its status.
// Includes random failures (e.g., running out of ingredients or the cook quitting mid-orderFinished).
func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++ // Increment the pizza orderFinished number

	if pizzaNumber <= numberOfPizzas {
		delay := rand.Intn(5) + 1 // Random delay for preparing the pizza
		fmt.Printf("Received orderFinished #%d \n", pizzaNumber)

		rnd := rand.Intn(12) // Random failure simulation

		msg := "Pizza #" + string(rune(pizzaNumber)) + " is ready"
		success := false // Default to failure

		if rnd < 5 {
			pizzasFailed++ // Increment failed pizzas
		} else {
			pizzasMade++ // Increment successful pizzas
		}

		total++ // Update the total orders

		fmt.Printf("Making #%d will take %d seconds\n", pizzaNumber, delay)

		// Simulate the delay in making a pizza
		time.Sleep(time.Duration(delay) * time.Second)

		// Conditions to determine success or failure
		if rnd <= 2 {
			msg = fmt.Sprint(" *** We ran out of ingredients for pizza #", pizzaNumber, " *** ")
		} else if rnd <= 4 {
			msg = fmt.Sprint(" *** The cook quit while making pizza #", pizzaNumber, " *** ")
		} else {
			success = true
			msg = fmt.Sprint("Pizza orderFinished #", pizzaNumber, " is ready")
		}

		return &PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}
	} else {
		// If no more pizzas can be made, return a message to indicate this
		return &PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     "No more pizzas",
			success:     false,
		}
	}
}

// main is the entry point of the program. It creates a producer, starts the pizzeria (producer),
// and processes the pizza orders (consumer).
func startMakingPizza() {
	// Seed the random number generator for random delays and outcomes
	rand.Seed(time.Now().UnixNano())

	// Print out a welcome message
	color.Cyan("The Pizzeria has started taking orders!")
	color.Cyan("---------------------------------------")

	// Create a producer
	pizzaJob := Producer{
		data: make(chan PizzaOrder), // Channel for pizza orders
		quit: make(chan chan error), // Channel for quit signals
	}

	// Run the producer in the background
	go pizzeria(&pizzaJob)

	// Consumer: Process orders from the producer
	for pizzaOrder := range pizzaJob.data {
		if pizzaOrder.pizzaNumber <= numberOfPizzas {
			// Handle successful or failed orders
			if pizzaOrder.success {
				color.Green(pizzaOrder.message)
				color.Green("Order #" + strconv.Itoa(pizzaOrder.pizzaNumber) + " is on the way!")
			} else {
				color.Red(pizzaOrder.message)
				color.Red("Order #" + strconv.Itoa(pizzaOrder.pizzaNumber) + " failed!")
			}
		} else {
			// No more pizzas to make
			color.Cyan("Done making pizzas!")
			err := pizzaJob.Close() // Close the producer
			if err != nil {
				color.Red("Error closing producer: %v", err)
			}
			break
		}
	}

	// Print summary at the end of the day
	color.Cyan("----------------")
	color.Cyan("Done for the day")

	color.Blue("Made: %d \n Failed: %d \n Total: %d \n", pizzasMade, pizzasFailed, total)

	// Evaluate the day's performance
	switch {
	case pizzasFailed > 9:
		color.Red("Awful day at the pizzeria!")
	case pizzasFailed >= 6:
		color.Yellow("Bad day at the pizzeria!")
	case pizzasFailed >= 4:
		color.HiYellow("Could be better at the pizzeria!")
	case pizzasFailed >= 2:
		color.HiGreen("Great day at the pizzeria!")
	case pizzasFailed >= 1:
		color.Green("Awesome day at the pizzeria!")
	default:
		color.Green("Best day at the pizzeria!")
	}
}
