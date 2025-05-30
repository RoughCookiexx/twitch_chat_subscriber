package twitch_chat_subscriber

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func SendRequestWithCallbackAndRegex(callbackURL string, regexPattern string) (string, error) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	config, err := loadConfig("config.json")
	if err != nil {
		log.Println("error loading config")
		return "", fmt.Errorf("error loading configuration: %w", err)
	}

	targetURL, err := url.Parse(config.TargetURL)
	if err != nil {
		log.Println("error parsing target URL from config")
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
		log.Println("error creating request")
		return "", fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error sending request")
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	log.Println("Sent subscriptions request to %s and received response code %s", callbackURL, resp.Status)
	// Return status code (e.g., "200 OK") and nil error if successful
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Println("status:", resp.Status)
		return resp.Status, nil
	}

	log.Println("Status: ", resp.Status)
	// Return status code and an error message for non-2xx responses
	return resp.Status, fmt.Errorf("request failed with status: %s", resp.Status)
}
