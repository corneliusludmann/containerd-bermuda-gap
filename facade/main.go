package main

import (
	"log"
)

func main() {
	log.Println("Hello, World!")
	reg, err := NewRegistry()
	if err != nil {
		log.Fatal(err)
	}
	reg.MustServe()
	log.Println("Bye!")
}
