package audit_test

import (
	"context"
	"testing"

	"github.com/Inengs/tweet-audit/src/archive"
	"github.com/Inengs/tweet-audit/src/audit"
	"github.com/Inengs/tweet-audit/src/config"
	"github.com/Inengs/tweet-audit/src/gemini"
	"golang.org/x/time/rate"
	"google.golang.org/api/googleapi"
)

func TestSplitIntoBatches(t *testing.T) {
	// create a slice of 45 tweets
	tweets := make([]archive.Tweet, 45)

	batches, err := audit.SplitIntoBatches(tweets, 20)

	if err != nil {
		t.Errorf("error: %v", err)
	}

	if len(batches[0]) != 20 {
		t.Errorf("error: the number of tweets in a batch is not equal to 20, %v", err)
	}

	if len(batches[2]) != 5 {
		t.Errorf("error: the number of the tweets in the last batch is not 5, %v", err)
	}

	if len(batches) != 3 {
		t.Errorf("the length of batches is not 3, %v", err)
	}
}

func TestAuditTweets_RetryOn429(t *testing.T) {
	callCount := 0
	

	mock := gemini.MockGeminiClient{
		AnalyzeTweetsFunc: func(tweets []archive.Tweet, username string, criteria config.Criteria) ([]gemini.FlaggedTweet, error) {
			callCount++
			if callCount < 3 {
				// simulate 429 for first 2 calls
				return nil, &googleapi.Error{Code: 429}
			}

			// succeed on 3rd Call
			return []gemini.FlaggedTweet{{TweetURL: "https://x.com/user/status/123"}}, nil
		},
	}

	tweets := make([]archive.Tweet, 5)
    ctx := context.Background()
    
	limiter := rate.NewLimiter(rate.Inf, 1)
    results, err := audit.AuditTweets(ctx, &mock, limiter, tweets, "user", config.Criteria{})
    
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    if callCount != 3 {
        t.Errorf("expected 3 calls, got %d", callCount)
    }
    if len(results) != 1 {
        t.Errorf("expected 1 flagged tweet, got %d", len(results))
    }
}

func TestAuditTweets_ContextCancellation(t *testing.T) {
	mock := &gemini.MockGeminiClient{
		// this is a placeholder mock for audit tweets
		AnalyzeTweetsFunc: func(tweets []archive.Tweet, username string, criteria config.Criteria) ([]gemini.FlaggedTweet, error) {
			return []gemini.FlaggedTweet{}, nil
		},
	}

	tweets := make([]archive.Tweet, 100) // 5 batches
	ctx, cancel := context.WithCancel(context.Background()) // this gives a context and a cancel
	cancel() // cancel immediately before the function starts

	limiter := rate.NewLimiter(rate.Inf, 1)
	results, err := audit.AuditTweets(ctx, mock, limiter, tweets, "user", config.Criteria{}) // at the start of each batch, the auditor will check ctx.Done(), since context has been cancelled, ctx.Done() fires immediately

	if err == nil {
		t.Errorf("expected context cancellation error, got nil") // we expect the cancellation error, which we had set off before
	}

	if results != nil {
		t.Errorf("expected no results after cancellation, got %d", len(results)) // we expect no results since, the error is set off on the first batch
	}
}

func TestAuditTweets_ReturnsFlaggedTweets(t *testing.T) {
	mock := &gemini.MockGeminiClient{
		AnalyzeTweetsFunc: func(tweets []archive.Tweet, username string, criteria config.Criteria) ([]gemini.FlaggedTweet, error) {
			return []gemini.FlaggedTweet{
				{TweetURL: "https://x.com/user/status/123", Deleted: false},
			}, nil
		},
	}

	tweets := make([]archive.Tweet, 5)
	ctx := context.Background()

	limiter := rate.NewLimiter(rate.Inf, 1) // unlimited rate for tests
	results, err := audit.AuditTweets(ctx, mock, limiter, tweets, "user", config.Criteria{})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 flagged tweet, got %d", len(results))
	}

	if results[0].TweetURL != "https://x.com/user/status/123" {
		t.Errorf("unexpected tweet URL: %s", results[0].TweetURL)
	}
}