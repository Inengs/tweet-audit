package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Inengs/tweet-audit/src/archive"
	"github.com/Inengs/tweet-audit/src/config"
	"google.golang.org/genai"
)

// this is the response shape of the tweets that are flagged
type FlaggedTweet struct {
	TweetURL string `json:"tweet_url"`
	Deleted bool `json:"deleted"`
	Reason string `json:"reason"`
	Confidence int `json:"confidence"`
}


// this is just a client shape for new client setup
type Client struct {
	model string
	client *genai.Client
}

// Create a new client
func NewClient(apiKey string) (*Client, error) {
	ctx := context.Background() // base context

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey, // define the api key in the client config
	}) // creates a new GenAI client
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &Client{ 
		model: "gemini-2.5-flash-lite", // use this version, seems like the best
		client: client,
	}, nil
}

// Analyze Tweets processes up to 20 tweets in one API call using batching
func (c *Client) AnalyzeTweets(tweets []archive.Tweet, username string, criteria config.Criteria) ([]FlaggedTweet, error) {
	if len(tweets) == 0 { // if the length of tweets less than 0
		return nil, nil
	}
	if len(tweets) > 20 { // if the length of tweets > 20
		return nil, fmt.Errorf("batch size exceeded (max 20)")
	}

	ctx := context.Background()

	prompt := BuildBatchPrompt(tweets, username, criteria)

	temp := float32(0.1)
	resp, err := c.client.Models.GenerateContent(
		ctx,
		c.model,
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			Temperature:     &temp,
			MaxOutputTokens: 1200,
		},
	)

	if err != nil {
		return nil, err
	}

	text := resp.Candidates[0].Content.Parts[0].Text
	return ParseFlaggedTweets(text)
} 

func BuildBatchPrompt(tweets []archive.Tweet, username string, criteria config.Criteria) string {
	var sb strings.Builder

	for i, t := range tweets {
		fmt.Fprintf(&sb, "Tweet %d | ID: %s | Text: \"%s\"\n\n", i+1, t.ID, t.Text)
	}

	return fmt.Sprintf(`You are auditing old tweets for deletion.

Rules:
- Forbidden Words: %s
- Forbidden Phrases: %s
- Unprofessional: %s
- Outdated Opinions: %s
- Custom Rules: %s
- Additional: %s

---

Analyze these tweets and return **ONLY** the ones that should be deleted.

Return a JSON array like this (empty array if none should be deleted):

[
  {
    "tweet_url": "https://x.com/%s/status/TWEET_ID_HERE",
    "deleted": false,
    "reason": "short reason",
    "confidence": 8
  }
]

IMPORTANT: Respond with ONLY a valid JSON array. No explanations, no text, no markdown. If no tweets should be flagged, respond with exactly: []

Tweets:
%s`, strings.Join(criteria.ForbiddenWords, ", "),
		strings.Join(criteria.ForbiddenPhrases, ", "),
		strings.Join(criteria.UnprofessionalPhrases, ", "),
		formatOutdatedOpinions(criteria.OutdatedOpinions),
		strings.Join(criteria.CustomRules, " | "),
		criteria.AdditionalInstructions,
		username,
		sb.String(),)
}

func ParseFlaggedTweets(responseText string) ([]FlaggedTweet, error) {
	// formats the result from gemini to make it more useful
	text := strings.TrimSpace(responseText)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var flagged []FlaggedTweet

	err := json.Unmarshal([]byte(text), &flagged) // converts the flagged tweets json into the required struct shape
	if err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response: %w", err)
	}

	return flagged, nil // returns the flagged tweets struct
}

func formatOutdatedOpinions(opinions []config.OutdatedOpinions) string {
	// base check
	if len(opinions) == 0 {
		return "None"
	}

	var sb strings.Builder
	for _, opinion := range opinions {
		fmt.Fprintf(&sb, "- %s: Used to believe \"%s\"", opinion.Topic, opinion.OldView)
		if opinion.NewView != "" {
			fmt.Fprintf(&sb, " → Now: \"%s\"", opinion.NewView)
		}
		if opinion.Since != "" {
			fmt.Fprintf(&sb, " (since %s)", opinion.Since)
		}
		sb.WriteString("\n")
	}

	return sb.String()
}