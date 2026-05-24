package output

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/Inengs/tweet-audit/src/gemini"
)

type CSVWriter struct {
    file   *os.File
    writer *csv.Writer
}


func CreateCSVFile(outputPath string) (*CSVWriter, error){
	file, err := os.Create(outputPath) // create the output file
	if err != nil {
		return nil, fmt.Errorf("error in creating file: %v", err)
	}
	writer := csv.NewWriter(file) // returns a new writer that writes to file

	// write a header row 
	writer.Write([]string{"tweet_url", "deleted"})

	// flush writes to disk
	writer.Flush()

	return &CSVWriter{
		file,
		writer,
	}, nil
}

func (c *CSVWriter)WriteFlaggedTweets(flaggedTweet *gemini.FlaggedTweet) error {
	flagged := []string{flaggedTweet.TweetURL, strconv.FormatBool(flaggedTweet.Deleted)} // convert the needed data to []string
	err := c.writer.Write(flagged) // write the data in memory
	if err != nil{
		// catch any error
		return fmt.Errorf("error from writing flagged tweet: %v", err)
	}

	// flush writes to file
	c.writer.Flush()

	return nil
}

func (c *CSVWriter) Close() error {
	c.writer.Flush()
	return c.file.Close()
}