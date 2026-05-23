package gemini_test

import (
	"strings"
	"testing"

	"github.com/Inengs/tweet-audit/src/archive"
	"github.com/Inengs/tweet-audit/src/config"
	"github.com/Inengs/tweet-audit/src/gemini"
)

func TestParseFlaggedTweets(t *testing.T) {
	jsonResponse := `[
  {
    "tweet_url": "https://x.com/myusername/status/123456",
    "deleted": true,
    "reason": "contains forbidden word",
    "confidence": 9
  }
]`

	flagged, err := gemini.ParseFlaggedTweets(jsonResponse)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(flagged) != 1  {
		t.Errorf("length of flagged tweet is more than one, error: %v", err)
	}

	if flagged[0].Deleted != true {
		t.Errorf("Deleted field isnt true, but it should be true")
	}

	if flagged[0].Reason != "contains forbidden word" {
		t.Errorf("the reason field does not match")
	}
}

func TestBuildBatchPrompt(t *testing.T) {
	criteria := config.Criteria{
		ForbiddenWords: []string{"retard", "faggot"},
		CustomRules: []string{"Avoid political rants"},
	}
	
	tweets := []archive.Tweet{
		{ID: "123456", Text: "This is a test tweet with bad word retard"},
		{ID: "789012", Text: "This is a normal clean tweet"},
	}

	prompt := gemini.BuildBatchPrompt(tweets, "myusername", criteria)

	if !strings.Contains(prompt, "Tweet 1 | ID: 123456") {
    	t.Errorf("prompt does not contain expected tweet format")
	}

	if !strings.Contains(prompt, "retard") {
		t.Errorf("prompt does not contain forbidden word")
	}

	if !strings.Contains(prompt, "myusername") {
		t.Errorf("prompt does not contain any username")
	}
}