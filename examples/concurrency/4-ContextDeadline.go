package main

import (
	"context"
	"fmt"
	"time"
)

func writeNums(ctx context.Context) {
	// WITH DEADLINE
	// deadline := time.Now().Add(1500 * time.Millisecond)
	// ctx, cancelCtx := context.WithDeadline(ctx, deadline)
	// ctx will close 1.5secs after deadline was generated
	// ----------------------------------------------------

	// WITH TIMEOUT
	ctx, cancelCtx := context.WithTimeout(ctx, 1500*time.Millisecond)
	// ctx will close 1.5secs after WithTimeout was called
	// -----------------------------------------------------

	defer cancelCtx() // not technically necessary in this function
	// useful if there were any return statements that could mean the later cancelCtx() was skipped

	printCh := make(chan int)
	go readNums(ctx, printCh)

	for num := 1; num <=3; num++ {
		select {
		case printCh <- num:
			time.Sleep(1 * time.Second) // here to demonstrate ctx deadline
		case <-ctx.Done():
			break 
		}
	} // don't send numbers to the channel after ctx is closed.
	// readNums is also checking the ctx 

	cancelCtx()
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("doSomething: finished:\n")
}

func readNums(ctx context.Context, printCh <- chan int) {
	for {
		select {
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				fmt.Printf("doAnother err: %s\n", err)
			}
			fmt.Printf("doAnother: finished\n")
			return
		case num := <-printCh:
			fmt.Printf("doAnother: %d\n", num)
		}
	}
}

func main() {
	ctx := context.Background()

	ctx = context.WithValue(ctx, "myKey", "value 1") 

	writeNums(ctx)
}