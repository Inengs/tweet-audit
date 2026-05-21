package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Inengs/tweet-audit/src/archive"
)

func main() {
	filepath := os.Args[1]

	tweets, err := archive.FileParser(filepath)
	
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, tweet := range tweets {
		fmt.Println(tweet)
	}
}