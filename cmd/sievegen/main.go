package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	// InitialValue is the starting value of the random walk.
	InitialValue = 100

	// Step is the range in which the walk can move.
	Step = 5

	// Interval is the frequency that data is generated.
	Interval = 100 * time.Millisecond
)

func main() {
	// Initialize the data to have a random value.
	var data struct {
		Value int `json:"value"`
	}
	data.Value = rand.Intn(InitialValue)

	// Perform a random walk on an interval.
	var encoder = json.NewEncoder(os.Stdout)
	for {
		// Print struct to stdout.
		if err := encoder.Encode(&data); err != nil {
			log.Fatalln("err:", err)
		}

		// Randomly walk the value.
		data.Value += rand.Intn(Step*2) - Step

		// Wait for a bit.
		time.Sleep(Interval)
	}
}
