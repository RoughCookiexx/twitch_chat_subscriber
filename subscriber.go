package twitch_chat_subscriber

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/RoughCookiexx/gg_twitch_types"
)

type StringResponse struct {
	Message string `json:"message"`
}

func SendRequestWithCallbackAndRegex(subscriptionURL string, callbackFunction func(twitch_types.Message)(string), regexPattern string, port int) (string, error) {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received callback message")

		var receivedMessage twitch_types.Message

		err := json.NewDecoder(r.Body).Decode(&receivedMessage)
		if err != nil {
			http.Error(w, "The JSON you sent was garbage.", http.StatusBadRequest)
			return
		}

		message := callbackFunction(receivedMessage)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(StringResponse{Message: message})
	})   

	targetURL, err := url.Parse(subscriptionURL)

	if err != nil {
		log.Println("Failed to parse target URL from string", subscriptionURL)
		return "", fmt.Errorf("Error parsing target URL: %s", subscriptionURL)
	}

	// Prepare query parameters
	queryParams := url.Values{}
	queryParams.Add("callbackURL", fmt.Sprintf("http://0.0.0.0:%d/callback", port))
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

	log.Printf("Sent subscriptions request and received response code %s", resp.Status)
	// Return status code (e.g., "200 OK") and nil error if successful
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Println("status:", resp.Status)
		return resp.Status, nil
	}

	log.Println("Status: ", resp.Status)
	// Return status code and an error message for non-2xx responses
	return resp.Status, fmt.Errorf("request failed with status: %s", resp.Status)
}
