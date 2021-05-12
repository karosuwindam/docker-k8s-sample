package main

import "fmt"

func main() {
	var ame280 Ame280
	ame280.Init()
	ame280.Test()
	fmt.Println(ame280.Read())
}
