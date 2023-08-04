package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	wg := sync.WaitGroup{}

	// time.Tick creates a channel and send the current time for every Tick.
	t := time.Tick(time.Second)

	// it will create 10 output channels. This channel will compete each others.
	dest := FanOut(t, 10)

	wg.Add(len(dest))
	for _, d := range dest {
		// creation of a local variable is important. If you use d directly, it will cause problems.
		// you have two option : 1- creation of a local variable 2-passing d as argument of the goroutine.
		i := d
		go func() {
			defer wg.Done()
			for val := range i {
				fmt.Println(val)
			}
		}()
	}

	wg.Wait()

}

// FanOut function takes input channel and number of destination channels will be created.
func FanOut(inputCh <-chan time.Time, n int) []<-chan time.Time {

	dest := make([]<-chan time.Time, 0)

	for i := 0; i < n; i++ {
		// creation of a dest channel
		d := make(chan time.Time, 0)

		dest = append(dest, d)

		go func() {
			// when inputCh is closed, defer will be executed and destination channels will be closed also.
			defer close(d)

			// each of this go routines will compete each other to get a value from input channel.
			for s := range inputCh {
				d <- s
			}
		}()
	}

	return dest
}
