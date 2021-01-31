package main

import "fmt"

func process(c chan int) {
	for i := 0; i <= 10; i++ {
		c <- i
	}
	c <- -1
}

func main() {
	channel := make(chan int)
	go process(channel)

	for {

		val := <-channel
		if val != -1 {
			fmt.Println(val)
		} else {
			close(channel)
			break
		}

	}
}
