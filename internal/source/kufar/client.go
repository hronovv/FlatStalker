package kufar

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) " +
	"AppleWebKit/537.36 (KHTML, like Gecko) " +
	"Chrome/120.0.0.0 Safari/537.36"

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) FetchAds(searchURL string) ([]Ad, error) {
	params, err := ParseSearchURL(searchURL)
	if err != nil {
		return nil, err
	}
	return c.FetchAdsByParams(params)
}

func (c *Client) FetchAdsByParams(params map[string]string) ([]Ad, error) {
	apiURL := BuildAPIURL(params)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("kufar API вернул %s: %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var payload struct {
		Ads []json.RawMessage `json:"ads"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	ads := make([]Ad, 0, len(payload.Ads))
	for _, raw := range payload.Ads {
		ad, err := AdFromAPI(raw)
		if err != nil {
			return nil, err
		}
		ads = append(ads, ad)
	}

	return ads, nil
}
