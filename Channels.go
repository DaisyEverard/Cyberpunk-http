package main

import ("fmt"
"sync")

func generateNumbers (total int, ch chan<- int, wg *sync.WaitGroup) {
	// The parameter type of chan<- int means the channel can only be written to, not read from
	// limit functionality unless needed
	defer wg.Done()


	for idx := 1; idx <= total; idx++{
	fmt.Printf("sending %d to channel\n", idx)
	ch <- idx
 }
}

func printNumbers(ch <-chan int,wg *sync.WaitGroup) {
	defer wg.Done()

	for num := range ch {
		fmt.Printf("read %d from channel\n", num)
	}
	// This loop reads from a channel until it closes
	// If you forgot to close the loop in main, the loop would never stop
	// and printNumbers could never finish.
}

func main() {
	var wg sync.WaitGroup
	numberChan := make(chan int) // makes a channel for ints

	wg.Add(2)
	go printNumbers(numberChan, &wg)
	generateNumbers(3, numberChan, &wg) // This is not a goroutine anymore.
	// main could close the channel before generateNumbers was finished
	// generate numbers try write to the channel and cause as panic type
	// "send on closed channel."
	
	close(numberChan) // closes the channel

	fmt.Println("Waiting for goroutime to finish...")
	wg.Wait()
	fmt.Println("Done!")
}