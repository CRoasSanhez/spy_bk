package core

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/revel/revel"
)

// HTTPMethod is a method for http request.
type HTTPMethod string

// Color type refers to colors
type AdviceType string

// Slack constants
const (
	// SlackWebHooksURL is url for web hooks services
	SlackWebHooksURL = "https://hooks.slack.com/services/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"

	// HTTP methods
	GET  HTTPMethod = "GET"
	POST HTTPMethod = "POST"

	// Advice types
	InfoEntry     AdviceType = "#e0e0eb" // Gray
	NewEntry      AdviceType = "#36a64f" // Green
	ErrEntry      AdviceType = "#ff0000" // Red
	ValidateEntry AdviceType = "#ff9900" // Orange

	CommonFooter = "Spychatter"
)

// SlackAttachment ...
type SlackAttachment struct {
	Fallback   string        `json:"fallback"`    // Notification text
	Color      string        `json:"color"`       // Color of attachment section
	Pretext    string        `json:"pretext"`     // Text after attachment
	AuthorName string        `json:"author_name"` // First text
	AuthorLink string        `json:"author_link"`
	Title      string        `json:"title"`
	Text       string        `json:"text"`
	Fields     []interface{} `json:"fields"`
	Footer     string        `json:"footer"`
	TimeStamp  int64         `json:"ts"`
}

// SlackMessage is a struct for send message to slack.
type SlackMessage struct {
	Text        string            `json:"-"`
	Attachments []SlackAttachment `json:"attachments"`
}

// Notify ...
func Notify(advice AdviceType, fallback, pretext, autor, autorURL, title, description string, fields []interface{}) error {

	switch revel.RunMode {
	case "candidate":
		pretext = "[ CANDIDATE ]  " + pretext
	case "dev":
		pretext = "[ DEV ]  " + pretext
	}

	return NotifyBySlack(SlackMessage{
		Text: "",
		Attachments: []SlackAttachment{
			{
				AuthorLink: autorURL,
				AuthorName: autor,
				Color:      string(advice),
				Fallback:   fallback,
				Fields:     fields,
				Footer:     CommonFooter,
				Pretext:    pretext,
				Text:       description,
				TimeStamp:  time.Now().Unix(),
				Title:      title,
			},
		},
	})
}

// NotifyBySlack sent message to slack
func NotifyBySlack(message SlackMessage) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := createHTTPRequest(POST, SlackWebHooksURL, data)
	if err != nil {
		return err
	}

	addHeaders(req, map[string]string{"Content-Type": "application/json"})

	if err = makeHTTPRequest(req); err != nil {
		return err
	}

	return nil
}

func createHTTPRequest(method HTTPMethod, url string, data []byte) (*http.Request, error) {
	return http.NewRequest(string(method), url, bytes.NewBuffer(data))
}

func addHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

func makeHTTPRequest(req *http.Request) error {
	// Create httpClient
	client := &http.Client{}

	// Send Request
	response, err := client.Do(req)
	if err != nil {
		return err
	}

	// Close when finish function
	defer response.Body.Close()

	return nil
}
