/*
You are given an API interface of some service for streaming data, which is the price of a “Ticker” (eg. BTC price from an exchange). The price of BTC can be different on different exchanges, so our target is to build a “fair” price for BTC, combining the price from different sources. Let’s say there are up to 100 possible exchanges from where the price can be streamed.
You need to develop an algorithm which uses these data streams as input and as output providing an online “index price” in the form of minute bars, where the bar is a pair (timestamp, price). Output can be provided in any form (file, console, etc.). An example output if the service is working for ~2 minutes would look like:
Timestamp, IndexPrice
1577836800, 100.1
1577836860, 102
Requirements:
Data from the streams can come with delays, but strictly in increasing time order for each stream. Stream can return an error, in that case the channel is closed. Bars timestamps should be solid minute as shown in example and provided in on-line manner, price should be the most relevant to the bar timestamp. How to combine different prices into the index is up to you.
Code should be written using Go language, apart from that you are free to choose how and what to do. You can also write some mock streams in order to test your code.
The interface is artificial, so if you need to change something or to have additional assumptions - you are free to do this, but don’t forget to mention that. Your code will be reviewed but won’t be executed on our side.
We expect source code to be published on GitHub and shared with us followed by the readme file with the description of the solution written in English. There might be a technical call after the task completion where we can discuss the solution in detail, ask some questions etc.
*/

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const startPrice = 100.0

type MockCh chan float64

func (m MockCh) start() {
	defer close(m)
	limit := 1 + rand.Intn(1000)
	price := startPrice
	for i := 0; i < limit; i++ {
		price += rand.NormFloat64() * float64(rand.Intn(11)-5) // change price
		m <- price
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	streams := make([]MockCh, 0, 3)
	for i := 0; i < cap(streams); i++ {
		streams = append(streams, make(chan float64, 10))
	}

	agg := make(chan float64, 100)
	defer close(agg)

	done := make(chan struct{})
	for _, s := range streams {
		go s.start()
	}

	// or we can send streams to queue
	for _, s := range streams {
		go func(c MockCh) {
			for msg := range c {
				agg <- msg
			}
			done <- struct{}{}
		}(s)
	}

	doneCount := 0
	for doneCount < len(streams) {
		lastTime := time.Now()
		select {
		case msg := <-agg:
			lastTime = checkTime(lastTime, msg)
			continue
		default:
		}
		select {
		case msg := <-agg:
			lastTime = checkTime(lastTime, msg)
		case <-done:
			doneCount++
		}
	}

}

// if val is last for minute - print
func checkTime(t time.Time, val float64) time.Time {
	newTime := time.Now()
	if newTime.Minute() == t.Minute() {
		rounded := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
		fmt.Printf("%d %.2f\n", rounded.Unix(), val)
		return newTime
	}
	return t
}
