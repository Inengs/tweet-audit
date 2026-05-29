package archive_test

import (
	"errors"
	"os"
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

func TestForInvalidPath(t *testing.T) {
	var filepath = "dummydata"

	_, err := archive.FileParser(filepath)

	if err == nil{
		t.Errorf("expected an error but got nil")
	}

	if !errors.Is(err, os.ErrNotExist) {
    	t.Errorf("expected ErrNotExist but got: %v", err)
	}
}

func TestForWrongFormat(t *testing.T) {
	var filepath = "testdata/invalid.json"

	_, err := archive.FileParser(filepath)

	if err == nil {
		t.Errorf("got a valid file: %v", err)
	}
}

func TestForEmptyArrayFile(t *testing.T) {
	var filepath = "testdata/empty_array.json"

	tweets, err := archive.FileParser(filepath)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	
	if len(tweets) != 0 {
		t.Errorf("expected empty slice, got %d tweets", len(tweets))
	}
}

func TestForEmptyFile(t *testing.T) {
	var filepath = "testdata/empty.json"

	_, err := archive.FileParser(filepath)

	if err == nil {
		t.Errorf("expected error: %v", err)
	}
}

