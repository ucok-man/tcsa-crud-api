package main

import (
	"fmt"
	"log"
)

func main() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", cfg)
}
