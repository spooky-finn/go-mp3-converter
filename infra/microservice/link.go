package linkservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

// Link -  is a microservice that is responsible for returning a downloadable url for a video
var ErrLinkEmptyResonse = errors.New("link request error: empty response")

type Result struct {
	Meta struct {
		Title  string `json:"title"`
		Source string `json:"source"`
	}
	Thumb string `json:"thumb"`
	URL   []struct {
		Url     string `json:"url"`
		Ext     string `json:"ext"`
		IsAudio bool   `json:"audio,omitempty"`
	} `json:"url"`
	DownloadURL string
}

func Fetch(url string) (*Result, error) {
	buf, err := json.Marshal(map[string]string{"url": url})
	if err != nil {
		panic(err)
	}

	client := http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Post(os.Getenv("LINK_URL"), "application/json", bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var results Result
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}

	downloadURL, err := parseResponse(results)
	if err != nil {
		return nil, err
	}

	results.DownloadURL = downloadURL
	return &results, nil
}

func parseResponse(data Result) (string, error) {
	// TODO: Add support audio link
	for _, result := range data.URL {
		if !result.IsAudio {
			return result.Url, nil
		}
	}
	return "", ErrLinkEmptyResonse
}
