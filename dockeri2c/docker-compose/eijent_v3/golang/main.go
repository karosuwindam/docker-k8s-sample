package main

import "context"

func Init() {

}

func main() {
	ctx := context.Background()
	<-ctx.Done()
}
