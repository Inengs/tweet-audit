package archive

import (
	"bytes"
	"encoding/json"
	"os"
)

type Tweet struct {
	ID        string `json:"id_str"`
	Text      string `json:"full_text"`
	CreatedAt string `json:"created_at"`
	Retweeted bool `json:"retweeted"`
}

type TweetWrapper struct {
	Tweet Tweet `json:"tweet"`
}

func FileParser(filepath string) ([]Tweet, error) {
	data, err := os.ReadFile(filepath) // read entire file and store as a byte slice
	if err != nil {
		return nil, err 
	}

	clean := bytes.TrimPrefix(data, []byte("window.YTD.tweets.part0 = ")) // trim out the first line

	var wrappers []TweetWrapper // used to go one level deeper to access the necessary fields
	err = json.Unmarshal(clean, &wrappers) // convert JSON to struct, so i can work with it
	if err != nil {
		return nil, err
	}

	var tweets []Tweet

	for _, w := range wrappers {
		tweets = append(tweets, Tweet{
			ID: w.Tweet.ID,
			Text: w.Tweet.Text,
			CreatedAt: w.Tweet.CreatedAt,
			Retweeted: w.Tweet.Retweeted,
		}) 
	}

	return tweets, nil
}