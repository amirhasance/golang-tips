package main

import (
	"fmt"
	"time"
)

func main() {

	dataChan := make(chan int, 100)

	go func() {
		defer func() {
			close(dataChan)
			fmt.Println("channel closed")
		}()
		for i := 0; i < 10; i++ {
			dataChan <- i
		}

	}()
	go func() {
		defer func() {
			fmt.Println("all data of channel is read")
		}()
		for d := range dataChan {
			time.Sleep(time.Second)
			fmt.Println("received", d)
		}
	}()
	time.Sleep(100 * time.Second)

}
