package main

import (
	"testing"
	"time"
)

func TestDine(t *testing.T) {
	eatTime = 0 * time.Second
	thinkTime = 0 * time.Second

	for i := 0; i < 10; i++ {
		orderFinished = []string{}
		dine()
		if len(orderFinished) != 5 {
			t.Error("Expected 5 orderFinished, got", len(orderFinished))
		}
	}
}

func TestDineWithVaryingDelay(t *testing.T) {
	var tests = []struct {
		name  string
		delay time.Duration
	}{
		{"0 second", 2 * time.Second},
		{"second", 200 * time.Millisecond},
		{"second", 500 * time.Millisecond},
	}

	for _, tt := range tests {
		orderFinished = []string{}
		eatTime = tt.delay
		thinkTime = tt.delay
		dine()
		if len(orderFinished) != 5 {
			t.Error("Expected 5 orderFinished, got", len(orderFinished))
		}
	}
}
