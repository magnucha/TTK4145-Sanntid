// Go 1.2
// go run helloworld_go.go

package main

import (
    . "fmt"
    //"runtime"
    "time"
)

func goroutine1(chan1 chan<- int, chan_finished chan<- int) {
    for j := 0; j < 1000000; j++ {
        chan1 <- 1;
    }
    chan_finished <- 1
}
func goroutine2(chan1 chan<- int, chan_finished chan<- int) {
    for k := 0; k < 1000000; k++ {
        chan1 <- -1;
    }
    chan_finished <- 1
}

func server(chan1 <-chan int, chan_finished <-chan int) int {
	var i int = 0
	var threads_finished int = 0
	for threads_finished < 2 {
		select {
			case msg1 := <- chan1:
				i += msg1
			case <- chan_finished:
				threads_finished++
		}
	}
	return i
}


func main() {
    //runtime.GOMAXPROCS(runtime.NumCPU())    // I guess this is a hint to what GOMAXPROCS does...
                                            // Try doing the exercise both with and without it!
    
    chan1 := make(chan int)
    chan2 := make(chan int)
    
    go goroutine1(chan1, chan2)                      // This spawns someGoroutine() as a goroutine
    go goroutine2(chan1, chan2)
    i := server(chan1, chan2)

    // We have no way to wait for the completion of a goroutine (without additional syncronization of some sort)
    // We'll come back to using channels in Exercise 2. For now: Sleep.
    time.Sleep(100*time.Millisecond)
    Println(i)
}
