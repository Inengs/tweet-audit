package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Inengs/tweet-audit/src/archive"
	"github.com/Inengs/tweet-audit/src/audit"
	"github.com/Inengs/tweet-audit/src/config"
	"github.com/Inengs/tweet-audit/src/gemini"
	"github.com/Inengs/tweet-audit/src/output"
)

func main() {
	c, err := config.LoadConfig("config.json")
	if err != nil{
		log.Fatalf("error in loading config files, %v", err)
	}

	fmt.Printf("Using API key: %s...\n", c.GeminiAPIkey[:10])
	
	archiveTweets, err := archive.FileParser(c.ArchivePath)
	if err != nil {
		log.Fatalf("failed to parse file, %v", err)
	}

	fmt.Printf("Parsed %d tweets\n", len(archiveTweets))

	geminiClient, err := gemini.NewClient(c.GeminiAPIkey)
	if err != nil {
		log.Fatalf("failed to create Gemini client: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	flaggedTweets, err := audit.AuditTweets(ctx, geminiClient, archiveTweets, c.Username, c.Criteria)
	if err != nil {
		log.Fatalf("failed to properly audit the tweets, %v", err)
	}

	fmt.Printf("Found %d flagged tweets\n", len(flaggedTweets))

	outputCSV, err := output.CreateCSVFile(c.OutputPath)
	if err != nil {
		log.Fatalf("failed to create csv file, %v", err)
	}

	for _, flaggedTweet := range flaggedTweets {
		err := outputCSV.WriteFlaggedTweets(&flaggedTweet)

		if err != nil {
			log.Fatalf("failed to write flagged tweets to csv file")
		}
	}

	fmt.Printf("Results saved to %s\n", c.OutputPath)


	err = outputCSV.Close()
	if err != nil {
		log.Fatalf("failed to close the file: %v", err )
	}
}