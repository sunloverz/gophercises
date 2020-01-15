package main

import (
	"fmt"
	"log"
	"os"
	"gophercises/link"
)

func main() {
	file, _ := os.Open("ex3.html")
	defer file.Close()
	links, err := link.Parse(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(links)
}
