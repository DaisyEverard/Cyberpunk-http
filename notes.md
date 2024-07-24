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

When part of a prgram writes to a channel, it will not do so again until another part reads from it. Same the other way round. Deadlocks happen when all blocking parts of the program (those waiting for something else to read or write) cannot be unblocked.



## Context

lets a function know about the environment it's being executed in
e.g. if a client disconnects it doesn't need process and send a response
