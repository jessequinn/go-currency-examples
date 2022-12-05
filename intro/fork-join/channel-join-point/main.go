package main

import (
	"fmt"
	"time"
)

type Example struct {
	x int
}

func main() {
	now := time.Now()
	done := make(chan Example)
	go func() {
		work()
		done <- Example{x: 64}
	}()

	c := <-done
	fmt.Println("elapsed:", time.Since(now), c.x)
	fmt.Println("done waiting, main exits")
}

func work() {
	time.Sleep(500 * time.Millisecond)
	fmt.Println("printing some stuff")
}
