package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	sources := make([]chan string, 20)

	// creation and population of source channels.
	for i := 0; i < 10; i++ {
		c := make(chan string)
		sources = append(sources, c)

		// goroutine below should get i value in this way. (IMPORTANT). otherwise it might be dangerous.
		chNum := i + 1

		// source channels will be populated with values from time.Tick channel. (trying to direction of values.)
		go func() {
			// different frequencies.
			t := time.Tick(time.Duration(chNum) * time.Second)
			for val := range t {
				str := fmt.Sprintf("%v (from: %d.source channel)", val, chNum)
				c <- str
			}
		}()
	}

	dest := FanIn(sources)

	for val := range dest {
		fmt.Println(val)
	}
}

func FanIn(sources []chan string) <-chan string {
	dest := make(chan string)

	wg := sync.WaitGroup{}

	for _, inputCh := range sources {

		c := inputCh
		wg.Add(1)

		go func() {
			defer wg.Done()

			for val := range c {
				dest <- val
			}

		}()
	}

	// wg.Wait is a blocking call. We need return dest channel, so wg.Wait will be in another goroutine.
	go func() {
		wg.Wait()
		close(dest)
	}()

	return dest
}
