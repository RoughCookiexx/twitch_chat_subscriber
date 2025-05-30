package twitch_chat_subscriber

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// Config struct to hold the target URL from config.json
type Config struct {
	TargetURL string `json:"targetURL"`
}

// loadConfig reads the targetURL from config.json
func loadConfig(filePath string) (Config, error) {
	var config Config

	configFile, err := os.Open(filePath)
	if err != nil {
		return config, fmt.Errorf("failed to open config file: %w", err)
	}
	defer configFile.Close()

	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal config JSON: %w", err)
	}

	if config.TargetURL == "" {
		return config, fmt.Errorf("targetURL not found or empty in config file")
	}

	return config, nil
}

// SendRequestWithCallbackAndRegex sends an HTTP GET request to a URL specified
// in config.json, appending callbackURL and regexPattern as query parameters.
// It returns the HTTP status code and an error message if any.
func SendRequestWithCallbackAndRegex(callbackURL string, regexPattern string) (string, error) {
	config, err := loadConfig("config.json")
	if err != nil {
		return "", fmt.Errorf("error loading configuration: %w", err)
	}

	targetURL, err := url.Parse(config.TargetURL)
	if err != nil {
		return "", fmt.Errorf("error parsing target URL from config: %w", err)
	}

	// Prepare query parameters
	queryParams := url.Values{}
	queryParams.Add("callbackURL", callbackURL)
	queryParams.Add("filterPattern", regexPattern)
	targetURL.RawQuery = queryParams.Encode()

	// Send HTTP GET request
	req, err := http.NewRequest("GET", targetURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Return status code (e.g., "200 OK") and nil error if successful
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return resp.Status, nil
	}

	// Return status code and an error message for non-2xx responses
	return resp.Status, fmt.Errorf("request failed with status: %s", resp.Status)
}
