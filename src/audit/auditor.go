package audit

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Inengs/tweet-audit/src/archive"
	"github.com/Inengs/tweet-audit/src/config"
	"github.com/Inengs/tweet-audit/src/gemini"
	"google.golang.org/api/googleapi"
)

// interface for mock gemini client
type TweetAnalyzer interface {
	AnalyzeTweets(tweets []archive.Tweet, username string, criteria config.Criteria) ([]gemini.FlaggedTweet, error)
}

func SplitIntoBatches(tweets []archive.Tweet, size int) ([][]archive.Tweet, error) {
	// guard against invalid batch sizes
	if size <= 0 {
		return nil, fmt.Errorf("tweet struct is empty")
	}

	var batches [][]archive.Tweet // array of batches of tweets

	for size < len(tweets) {
		tweets, batches = tweets[size:], append(batches, tweets[0:size:size]) // loop through 
	}

	// only append the remainder if there are actually left
	if len(tweets) > 0 {
		batches = append(batches, tweets)
	}

	return batches, nil
}

func AuditTweets(
	ctx context.Context,
	client TweetAnalyzer,
	tweets []archive.Tweet,
	username string,
	criteria config.Criteria,
) ([]gemini.FlaggedTweet, error) {
	batches, err := SplitIntoBatches(tweets, 20) // divide tweets into batches
	if err != nil {
		return nil, fmt.Errorf("failed to split the tweets into batches, error: %v", err)
	}

	var allResults []gemini.FlaggedTweet

	for _, batch := range batches { 
		// for each batch

		select {
    	case <-ctx.Done():
			// stops the flow safely in the case of any cancellation signal
        	return nil, ctx.Err()
    	default:
		}

		for retry := 0; retry < 3; retry++ {
			// for each retry

			results, err := client.AnalyzeTweets(batch, username, criteria) // call analyze tweets

			if err != nil {
				// if there is an error
				var apiErr *googleapi.Error
				if errors.As(err, &apiErr) && apiErr.Code == 429 { // if error is rate limit
					select {
					case <-ctx.Done():
						return nil, ctx.Err()
					case <-time.After(5 * time.Second):
					}
					continue 
				} else {
					// if not 429 error. log and move onto the next cycle of batches
					fmt.Printf("error is not a rate limiting error, err: %v", err)
					break
				}
			}

			// if no error, append results
			allResults = append(allResults, results...)

			// break out of retry loop after it has been successfully appended
			break
		}
	}

	return allResults, nil
}