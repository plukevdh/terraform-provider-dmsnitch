package dmsnitch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go/types"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// Code borrowed heavily from the Bitbucket provider.

type Error struct {
	Validations types.Array `json:"validations,omitempty"`
	Message     string      `json:"error,omitempty"`
	Type        string      `json:"type,omitempty"`
	StatusCode  int
	Endpoint    string
}

func (e Error) Error() string {
	return fmt.Sprintf("API Error: %d %s %s", e.StatusCode, BaseURL+e.Endpoint, e.Message)
}

const (
	BaseURL string = "https://api.deadmanssnitch.com/v1/"
)

type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

func (c *Client) Do(method, endpoint string, payload *bytes.Buffer) (*http.Response, error) {
	apiEndpoint := BaseURL + endpoint
	log.Printf("[DEBUG] Sending request to %s %s", method, apiEndpoint)

	var bodyreader io.Reader

	if payload != nil {
		log.Printf("[DEBUG] With payload %s", payload.String())
		bodyreader = payload
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, method, apiEndpoint, bodyreader)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.APIKey, ":")
	req.Header.Add("Content-Type", "application/json")
	req.Close = true

	resp, err := c.HTTPClient.Do(req)
	log.Printf("[DEBUG] Resp: %v Err: %v", resp, err)
	if resp.StatusCode >= 400 || resp.StatusCode < 200 {
		apiError := Error{
			StatusCode: resp.StatusCode,
			Endpoint:   endpoint,
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Resp Body: %s", string(body))

		err = json.Unmarshal(body, &apiError)
		if err != nil {
			apiError.Message = string(body)
		}

		return resp, error(apiError)
	}
	return resp, err
}

func (c *Client) Get(endpoint string) (*http.Response, error) {
	return c.Do("GET", endpoint, nil)
}

func (c *Client) Post(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
	return c.Do("POST", endpoint, jsonpayload)
}

func (c *Client) Patch(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
	return c.Do("PATCH", endpoint, jsonpayload)
}

func (c *Client) Delete(endpoint string) (*http.Response, error) {
	return c.Do("DELETE", endpoint, nil)
}
