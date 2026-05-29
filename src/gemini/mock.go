package gemini

import (
	"github.com/Inengs/tweet-audit/src/archive"
	"github.com/Inengs/tweet-audit/src/config"
)

type MockGeminiClient struct {
	AnalyzeTweetsFunc func(tweets []archive.Tweet, username string, criteria config.Criteria) ([]FlaggedTweet, error)
}

// AnalyzeTweets implements the GeminiClient interface
func (m *MockGeminiClient) AnalyzeTweets(tweets []archive.Tweet, username string, criteria config.Criteria) ([]FlaggedTweet, error) {
	if m.AnalyzeTweetsFunc != nil {
		return m.AnalyzeTweetsFunc(tweets, username, criteria)
	}
	return nil, nil // default behavior
}