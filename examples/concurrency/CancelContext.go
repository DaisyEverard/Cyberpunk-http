package main

import (
	"context"
	"fmt"
	"time"
)

func doSomething(ctx context.Context) {
	ctx, cancelCtx := context.WithCancel(ctx)

	printCh := make(chan int)
	go doAnother(ctx, printCh)

	for num := 1; num <=3; num++ {
		printCh <- num
	}
	cancelCtx()

	time.Sleep(100 * time.Millisecond)
	fmt.Printf("doSomething: finished:\n")
}

func doAnother(ctx context.Context, printCh <- chan int) {
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
    // While the behaviour is always that 3 numbers are printed, i believe this is a race condition.
	// The context is cancelled immediatedly after the 3rd number is written
	// This could mean the select will run the return statement before reading the 3rd number.

	// You could introduce another channel to output to every time the goroutine has processed the number
	// Check in the main function to check that 3 things have been processed (But this isn't a common thing to do,
	// The reason to cancel context is because you want it to stop processing)
}

func main() {
	ctx := context.Background() // create a new empty context

	ctx = context.WithValue(ctx, "myKey", "value 1") // return new context with 1 added value

	doSomething(ctx)
}