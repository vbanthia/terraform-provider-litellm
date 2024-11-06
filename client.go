package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	APIBase    string
	APIKey     string
	httpClient *http.Client
}

func NewClient(apiBase, apiKey string) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &Client{
		APIBase:    apiBase,
		APIKey:     apiKey,
		httpClient: &http.Client{Transport: tr},
	}
}

// Team-related methods
func (c *Client) CreateTeam(team map[string]interface{}) (map[string]interface{}, error) {
	return c.sendRequest("POST", "/team/new", team)
}

func (c *Client) GetTeam(teamID string) (map[string]interface{}, error) {
	return c.sendRequest("GET", fmt.Sprintf("/team/info?team_id=%s", teamID), nil)
}

func (c *Client) UpdateTeam(team map[string]interface{}) (map[string]interface{}, error) {
	return c.sendRequest("POST", "/team/update", team)
}

func (c *Client) DeleteTeam(teamID string) error {
	payload := map[string]interface{}{
		"team_ids": []string{teamID},
	}
	_, err := c.sendRequest("POST", "/team/delete", payload)
	return err
}

func (c *Client) sendRequest(method, path string, body interface{}) (map[string]interface{}, error) {
	url := c.APIBase + path

	var req *http.Request
	var err error

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %v", err)
		}
		log.Printf("Making %s request to %s with body:\n%s", method, url, string(jsonBody))
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	} else {
		log.Printf("Making %s request to %s", method, url)
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	log.Printf("Response status: %d", resp.StatusCode)
	log.Printf("Response body: %s", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		if method == "POST" && (len(bodyBytes) == 0 || string(bodyBytes) == "null") {
			return make(map[string]interface{}), nil
		}
		return nil, fmt.Errorf("error parsing response JSON: %v\nResponse body: %s", err, string(bodyBytes))
	}

	return result, nil
}
