package testing

import (
	"fmt"
)

//noinspection GoUnusedFunction
func talk() {
	name := "john"
	fmt.Println("Hello there", name)
	greeting("Hello to you too.")
}

func greeting(say string) {
	fmt.Println(say)
}
