package output_test

import (
	"encoding/csv"
	"os"
	"testing"

	"github.com/Inengs/tweet-audit/src/gemini"
	"github.com/Inengs/tweet-audit/src/output"
)

func TestCreateCSVFile(t *testing.T) {
	csvWriter, err := output.CreateCSVFile("test.csv")
	defer os.Remove("test.csv") // cleanup after test so that the file doesnt remain after test

	if err != nil {
		t.Errorf("error creating file: %v", err)
	}

	_, err = os.Stat("test.csv")
	if err != nil {
		t.Errorf("file was not created, %v", err)
	}

	csvWriter.Close()
}

func TestWriteFlaggedTweets(t *testing.T) {
	csvWriter, _ := output.CreateCSVFile("test.csv") // create csv file
    defer os.Remove("test.csv") // delete the file after test

	tweet := &gemini.FlaggedTweet{TweetURL: "https://x.com/user/status/123", Deleted: false}

	err := csvWriter.WriteFlaggedTweets(tweet)
	if err != nil {
		t.Errorf("failed to write to CSV, %v", err)
	}

	file, err := os.Open("test.csv") // open file
	if err != nil {
		t.Errorf("failed to open csv file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file) // create a reader

	csvWriter.Close() // close the file after writing to it

	allRows, err := reader.ReadAll() // read all the rows
	if err != nil {
		t.Errorf("failed to read from the csv: %v", err)
	}

	if len(allRows) != 2{
		t.Errorf("failed to write to the csv: %v", err)
	}
}