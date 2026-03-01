package ghost

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client talks to the Ghost Admin API.
type Client struct {
	BaseURL     string
	AdminAPIKey string
	HTTPClient  *http.Client
}

// GhostError represents an error response from the Ghost API.
type GhostError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Context string `json:"context,omitempty"`
}

// GhostErrorResponse wraps the errors array Ghost returns.
type GhostErrorResponse struct {
	Errors []GhostError `json:"errors"`
}

func (e *GhostErrorResponse) Error() string {
	if len(e.Errors) == 0 {
		return "unknown Ghost API error"
	}
	msg := e.Errors[0].Message
	if e.Errors[0].Context != "" {
		msg += " — " + e.Errors[0].Context
	}
	return msg
}

// NewClient creates a Ghost API client.
func NewClient(baseURL, adminAPIKey string) *Client {
	// normalize URL — strip trailing slash
	baseURL = strings.TrimRight(baseURL, "/")

	return &Client{
		BaseURL:     baseURL,
		AdminAPIKey: adminAPIKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest performs an authenticated request to the Ghost Admin API.
func (c *Client) doRequest(method, path string, body io.Reader, contentType string) ([]byte, error) {
	token, err := GenerateJWT(c.AdminAPIKey)
	if err != nil {
		return nil, fmt.Errorf("generating auth token: %w", err)
	}

	fullURL := c.BaseURL + "/ghost/api/admin" + path

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	req.Header.Set("Authorization", "Ghost " + token)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var ghostErr GhostErrorResponse
		if json.Unmarshal(respBody, &ghostErr) == nil && len(ghostErr.Errors) > 0 {
			return nil, &ghostErr
		}
		return nil, fmt.Errorf("Ghost API error (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Get performs a GET request.
func (c *Client) Get(path string, params url.Values) ([]byte, error) {
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	return c.doRequest(http.MethodGet, path, nil, "")
}

// Post performs a POST request with JSON body.
func (c *Client) Post(path string, body io.Reader) ([]byte, error) {
	return c.doRequest(http.MethodPost, path, body, "application/json")
}

// Put performs a PUT request with JSON body.
func (c *Client) Put(path string, body io.Reader) ([]byte, error) {
	return c.doRequest(http.MethodPut, path, body, "application/json")
}

// Delete performs a DELETE request.
func (c *Client) Delete(path string) ([]byte, error) {
	return c.doRequest(http.MethodDelete, path, nil, "")
}

// PostMultipart performs a POST with a custom content type (for image uploads).
func (c *Client) PostMultipart(path string, body io.Reader, contentType string) ([]byte, error) {
	return c.doRequest(http.MethodPost, path, body, contentType)
}

// Site fetches site information. Doubles as an auth check.
type SiteInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Logo        string `json:"logo,omitempty"`
	URL         string `json:"url"`
	Version     string `json:"version"`
}

type siteResponse struct {
	Site SiteInfo `json:"site"`
}

func (c *Client) GetSite() (*SiteInfo, error) {
	data, err := c.Get("/site/", nil)
	if err != nil {
		return nil, err
	}

	var resp siteResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing site response: %w", err)
	}

	return &resp.Site, nil
}
