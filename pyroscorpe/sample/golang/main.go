package main

import (
	"context"
	"fmt"
	"pyro-sample/pyroscopesetup"
	"time"
)

func loop1(ch chan<- int, ctx context.Context) {
	fmt.Println("loop1 Start")
	count := 0
loop1:
	for {
		select {
		case <-ctx.Done():
			break loop1
		case <-time.After(time.Microsecond * 500):
			count++
			// fmt.Println("input:", count)
			ch <- count
		}
	}
	fmt.Println("loop1 End")
}

func loop2(ch <-chan int, ctx context.Context) {
	fmt.Println("loop1 Start")
loop2:
	for {
		select {
		case <-ctx.Done():
			break loop2
		case count := <-ch:
			for i := 0; i < count; i++ {
			}
			fmt.Println("output:", count)
		}
	}
	fmt.Println("loop1 End")

}

func main() {
	py := pyroscopesetup.Setup()
	py.Run()

	ch1 := make(chan int, 20)
	ctx, cannel := context.WithCancel(context.Background())
	go loop1(ch1, ctx)
	go loop2(ch1, ctx)
	<-ctx.Done()
	close(ch1)
	cannel()
	time.Sleep(time.Microsecond * 500)
}
