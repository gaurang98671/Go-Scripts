package main

import (
	"fmt"
	"sync"
)

func process(wg *sync.WaitGroup) {
	fmt.Println("Executing go routine")
	wg.Done()
}
func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go process(&wg)
	go process(&wg)

	wg.Wait()
	fmt.Println("Done executing")

}
