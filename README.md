# Tweet audit platform

A CLI tool that audits your Twitter/X archive using Gemini AI and flags tweets for deletion based on your own criteria

## Prerequisites

- Go 1.21+
- A Twitter/X archive
- A Gemini API key

## Installation

```bash
git clone https://github.com/Inengs/tweet-audit
cd tweet-audit
go mod download
```

## Getting Your X Archive

1. Go to X.com → Settings → Your Account → Download an archive of your data
2. Request the archive and wait for the email
3. Download and extract the ZIP file
4. Your tweet data is at `data/tweets.js` inside the extracted folder

## Getting a Gemini API Key

1. Go to [https://aistudio.google.com](https://aistudio.google.com)
2. Sign in with your Google account
3. Click "Get API Key" and create a new key
4. Copy the key — you'll need it for your config

## Configuration

Copy the example config and fill in your details:

```bash
cp config.example.json config.json
```

Edit `config.json`:

```json
{
  "gemini_api_key": "YOUR_KEY_HERE",
  "username": "your_x_username",
  "archive_path": "/path/to/data/tweets.js",
  "output_path": "./flagged_tweets.csv",
  "criteria": {
    "forbidden_words": ["example"],
    "forbidden_phrases": ["example phrase"],
    "unprofessional_phrases": [],
    "outdated_opinions": [
      {
        "topic": "crypto",
        "old_view": "crypto is the future",
        "new_view": "I no longer hold this view",
        "since": "2024"
      }
    ],
    "custom_rules": [
      "flag any tweet that could embarrass me in a job interview"
    ],
    "additional_instructions": ""
  }
}
```

## Running the Tool

```bash
go run main.go
```

## Example Output

The tool writes a CSV file to your configured `output_path`:
tweet_url,deleted
https://x.com/username/status/123456,false
https://x.com/username/status/789012,false

Open the CSV, review each flagged tweet, and delete manually on X.

## Notes

- Press Ctrl+C to stop the audit safely — progress is saved incrementally
