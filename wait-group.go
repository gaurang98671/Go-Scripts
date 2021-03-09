package main

import (
	"fmt"
	"sync"
	"time"
)

func process(wg *sync.WaitGroup) {
	fmt.Println("Executing go routine")
	wg.Done()
}

func process2(wg *sync.WaitGroup) {
	time.Sleep(time.Duration(1 * time.Second))
	fmt.Println("Executing process 2")
	wg.Done()

}
func main() {
	var wg sync.WaitGroup
	wg.Add(4)

	go process2(&wg)
	go process2(&wg)
	go process(&wg)
	go process(&wg)

	wg.Wait()
	fmt.Println("Done executing")

}
