package output

import (
	"fmt"
)

func talk() {
	name := "john snow"
	fmt.Println("Hello there", name)
	greeting("Hello to you too.")
}

func greeting(say string) {
	fmt.Println(say)
}
