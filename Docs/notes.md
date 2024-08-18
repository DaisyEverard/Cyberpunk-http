## Concurrency

You need to make sure the main function can't finish before other functions finish running
see Waitgroup.go

goroutines can't have return values as the values with be thrown out once the function completes. Instead, you have to use channels.
Channels can only send one datatype

`make(chan int)` - makes a channel for integers
`intChan <- 10` - channel followed by `<-` means writing to the channel
`intVar := <- intChan` - `<-` followed by channel meands reading from the channel
`func readChannel(ch <- chan int) {}` read-only channel
`func writeChannel(ch chan <- int) {}` write only channel
`close(channelName)` - close a channel

```
for num := range intChan {
    if num < 1 {
        break
    }
 }
```

When reading a channel with the range keyword, you can set up a loop which reads the next value from the channel each time. It will read until the channel is closed or the loop is exited in another way, such as break.

When part of a porgram writes to a channel, it will not do so again until another part reads from it. Same the other way round. Deadlocks happen when all blocking parts of the program (those waiting for something else to read or write) cannot be unblocked.

## Context

lets a function know about the environment it's being executed in
e.g. if a client disconnects it doesn't need to process and send a response

`ctx context.Context` - common variable name ctx for context. Type from context package

There are 2 ways to create an empty/starting context.
Both have the same functionality, but signal different uses
`context.TODO()` - placeholder, unsure which context to use
`context.Background()` - start a known context.
`ctx.Value("key")` - access values

`context.WithValue` -3 params. Parent context.Context, key, value.
return value is a new context.Context.
Values are immutable
If you create a new context and use the same key that exists on the parent context, the old value of the key will be overwitten in the new context. This is acheved by wrapping the old context in a new one.
In the following example, ctx1's key will have value 1, but ctx2 will have value 2

```
ctx1 := context.WithValue(ctx, "a key", "value 1")
ctx2 := context.WithValue(ctx1, "a key", "value 2")
```

Context have a Done channel which is closed and start returning nil once the context is closed. Closing a context will close all child contexts.

Select statment allows program to try reading from or writing to a number of channels at the same time. Once channel operation per select statement.
When in a loop, program can do a number of channel operations whne one becomes available.

```
ctx := context.Background()
resultsCh := make(chan *WorkResult)

for {
	select {
	case <- ctx.Done():
		// The context is over, stop processing results
		return
	case result := <- resultsCh:
		// Process the results received from somewhere else
	}
}
```

There are 2 case statement. when the select statement is run, it will watch the case stement until ONE can be executed.
ctx.Done() does not have anything to read until the channel is closed
If the channel is still open, the first thing availiable to run will be reading from resutsCh.
If there are multiple cases which could be run, there is no guarantee of which will be run first

There are multiple ways to end a context.
Cancelling is the most controllable
`ctx, cancelCtx := context.WithCancel(ctx)` - withCancel returns a context and cancel method

`ctx, cancelContext := ctx.WithDeadline(context.Context, time.Time)`
This context will automatically close after time.Time

`ctx.WithTimeout(context.Context, time.Duration)`
This context will automatically close after the time.Duration amount has elapsed.
