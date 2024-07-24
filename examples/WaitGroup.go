package main

import ("fmt"
"sync")

func generateNumbers (total int, wg *sync.WaitGroup) {
	defer wg.Done() // the method will use defer to call Done once it's finished running and reduce the waitgroup count by 1
	for idx := 1; idx <= total; idx++{
	fmt.Printf("Generating number %d\n", idx)
 }
}

func printNumbers(wg *sync.WaitGroup) {
	defer wg.Done()
	for idx := 1; idx <= 3; idx++ {
		fmt.Printf("Printing number %d\n", idx)
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(2) // This tells the wait group to wait for 2 functions to finish
	go printNumbers(&wg) // putting the go keyword in from of the function call makes it run as a goroutine
	go generateNumbers(3, &wg)

	fmt.Println("Waiting for goroutime to finish...")
	wg.Wait() // This step will not be passed until the wait group count reaches 0
	fmt.Println("Done!")
}