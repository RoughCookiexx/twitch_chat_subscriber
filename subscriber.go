package twitch_chat_subscriber

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

func SendRequestWithCallbackAndRegex(subscriptionURL string, callbackURL string, regexPattern string) (string, error) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	targetURL, err := url.Parse(subscriptionURL)

	if err != nil {
		log.Println("Failed to parse target URL from string", subscriptionURL)
		return "", fmt.Errorf("Error parsing target URL: %s", subscriptionURL)
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
