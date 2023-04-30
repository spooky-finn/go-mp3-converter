package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"3205.team/go-mp3-converter/cfg"
)

// Link -  is a microservice that is responsible for returning a downloadable url for a video
var ErrLinkEmptyResonse = errors.New("link request error: empty response")

type Result struct {
	Meta struct {
		Title  string `json:"title"`
		Source bool   `json:"source"`
	}
	Thumb string `json:"thumb"`
	URL   []struct {
		Url     string `json:"url"`
		Ext     string `json:"ext"`
		IsAudio bool   `json:"audio,omitempty"`
	} `json:"url"`
	DownloadURL string
}

func LinkRequest(url string) (*Result, error) {
	buf, err := json.Marshal(map[string]string{"url": url})
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(cfg.AppConfig.LinkURL, "application/json", bytes.NewReader(buf))
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
