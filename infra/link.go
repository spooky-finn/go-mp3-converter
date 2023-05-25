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

type LinkResultData struct {
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

func Fetch(url string) (*LinkResultData, error) {
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

	var linkData LinkResultData
	err = json.Unmarshal(body, &linkData)
	if err != nil {
		return nil, err
	}

	downloadURL, err := parseResponse(linkData)
	if err != nil {
		return nil, err
	}

	linkData.DownloadURL = downloadURL
	return &linkData, nil
}

func parseResponse(data LinkResultData) (string, error) {
	// TODO: Add support audio link
	for _, result := range data.URL {
		if !result.IsAudio {
			return result.Url, nil
		}
	}
	return "", ErrLinkEmptyResonse
}
