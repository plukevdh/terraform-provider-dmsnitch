package dmsnitch

import (
	"fmt"
	"net/http"
	"bytes"
	"log"
	"io"
	"io/ioutil"
	"encoding/json"
	"go/types"
)

// Code borrowed heavily from the Bitbucket provider.

type Error struct {
	Validations types.Array `json:"validations,omitempty"`
	Message    string `json:"error,omitempty"`
	Type       string `json:"type,omitempty"`
	StatusCode int
	Endpoint   string
}

func (e Error) Error() string {
	return fmt.Sprintf("API Error: %d %s %s", e.StatusCode, e.Endpoint, e.Message)
}

const (
	BaseUrl string = "https://api.bitbucket.org/v1/"
)

type DMSnitchClient struct {
	ApiKey   string
	HTTPClient *http.Client
}

func (c *DMSnitchClient) Do(method, endpoint string, payload *bytes.Buffer) (*http.Response, error) {

	api_endpoint := BaseUrl + endpoint
	log.Printf("[DEBUG] Sending request to %s %s", method, api_endpoint)

	var bodyreader io.Reader

	if payload != nil {
		log.Printf("[DEBUG] With payload %s", payload.String())
		bodyreader = payload
	}

	req, err := http.NewRequest(method, api_endpoint, bodyreader)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.ApiKey, "")
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

func (c *DMSnitchClient) Get(endpoint string) (*http.Response, error) {
	return c.Do("GET", endpoint, nil)
}

func (c *DMSnitchClient) Post(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
	return c.Do("POST", endpoint, jsonpayload)
}

func (c *DMSnitchClient) Patch(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
	return c.Do("PATCH", endpoint, jsonpayload)
}

func (c *DMSnitchClient) Delete(endpoint string) (*http.Response, error) {
	return c.Do("DELETE", endpoint, nil)
}
