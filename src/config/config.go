package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Filtering criteria
type Criteria struct {
	// Simple keywords
	ForbiddenWords []string `json:"forbidden_words"`
	ForbiddenPhrases []string `json:"forbidden_phrases"`
	UnprofessionalPhrases []string `json:"unprofessional_phrases"`

	// custom rules
	CustomRules []string `json:"custom_rules"`

	// professional check
	ProfessionalCheck bool `json:"json:professional_check"`

	// outdated opinions
	OutdatedOpinions []OutdatedOpinions `json:"outdated_opinions"`

	// Context for the AI
	AdditionalInstructions string `json:"additional_instructions"`
}

type OutdatedOpinions struct {
	Topic string `json:"topic"`
	OldView string `json:"old_view"`
	NewView string `json:"new_view"`
	Since string `json:"since"`
}


// struct for config.json
type Config struct {
	GeminiAPIkey string `json:"gemini_api_key"`
	Username     string `json:"username"`
	ArchivePath  string `json:"archive_path"`
	OutputPath   string `json:"output_path"`
	Criteria     Criteria `json:"criteria"`
}

// To load config details
func LoadConfig(path string) (*Config, error) {
	configData, err := os.ReadFile(path) // read config.json
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config.json not found; copy config.example.json to config.json and fill in your details")
		}

		if 	os.IsPermission(err) {
			return nil, fmt.Errorf("permission denied reading config.json; check the file permissions")
		}

		return nil, fmt.Errorf("could not read config file at %s: %w", path, err)
	}

	var config Config
	err = json.Unmarshal(configData, &config) // collect raw json from config.json and convert to config struct
	if err != nil {
		return nil, fmt.Errorf("invalid JSON in config.json: %w", err)
	}

	// check for Gemini API key
	if config.GeminiAPIkey == "" {
		return nil, fmt.Errorf("gemini_api_key is required in config.json")
	}

	// check for Archive path
	if config.ArchivePath == "" {
		return nil, fmt.Errorf("archive path is required in config.json")
	}

	if config.OutputPath == "" {
		return nil, fmt.Errorf("output Path are required in the config.json")
	}

	if config.Username == "" {
		return nil, fmt.Errorf("your twitter username is required in config.json")
	}

	return &config, nil
}