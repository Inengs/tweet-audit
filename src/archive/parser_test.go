package archive_test

import (
	"testing"

	"github.com/Inengs/tweet-audit/src/archive"
)

func TestFileParser(t *testing.T) {
	var filepath = "testdata/tweets.json"

	tweets, err := archive.FileParser(filepath)

	if err != nil {
		t.Errorf("failed to parse file: %v", err)
	}

	if len(tweets) != 2 {
		t.Errorf("expected 2 tweets, got %d", len(tweets))
	}
}