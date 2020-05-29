{{ reserveImport "fmt"  }}

func talk() {
    name := "{{.Name}}"
    fmt.Println("Hello there", name)
    greeting("{{$.Greeting}}")
}

func greeting(say string) {
    fmt.Println(say)
}