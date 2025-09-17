/*
   Copyright (C) 2025 Rodolfo González González.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.

*/

package whoisparserai

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	resty "resty.dev/v3"
)

//-----------------------------------------------------------------------------

// Request structure for Azure OpenAI API
type ChatCompletionRequest struct {
	Messages    []Message `json:"messages"`
	MaxTokens   *int      `json:"max_tokens,omitempty"`
	Temperature *float64  `json:"temperature,omitempty"`
	TopP        *float64  `json:"top_p,omitempty"`
	Stop        []string  `json:"stop,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response structure from Azure OpenAI API
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Error response structure
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

type AzureOpenAIClient struct {
	APIKey     string
	Endpoint   string
	Deployment string
	Model      string
	APIVersion string
}

func NewAzureOpenAIClient(apiKey string, endpoint string, model string) *AzureOpenAIClient {
	return &AzureOpenAIClient{
		APIKey:     apiKey,
		Endpoint:   endpoint,
		Model:      model,
		APIVersion: "2025-01-01-preview",
	}
}

//-----------------------------------------------------------------------------

// Helper functions to create pointers
func intPtr(i int) *int             { return &i }
func float64Ptr(f float64) *float64 { return &f }

//-----------------------------------------------------------------------------

func (c *AzureOpenAIClient) ChatCompletion(req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Construct the URL
	url := fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s",
		c.Endpoint, c.Model, c.APIVersion)

	// Marshal request to JSON
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Create HTTP request
	client := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(8 * time.Second).
		AddRetryConditions(
			func(r *resty.Response, err error) bool {
				return err != nil || r.StatusCode() >= 500
			})

	// Request
	resp, err := client.R().
		SetBody(data).
		SetHeader("Content-Type", "application/json").
		SetHeader("api-key", c.APIKey).
		Post(url)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &chatResp, nil
}

//-----------------------------------------------------------------------------

func (c *AzureOpenAIClient) ParseWhois(whoisText string) (map[string]interface{}, error) {
	roleSystem := `
You are a parser that extracts WHOIS data into JSON. Extract the following data
from the given whois output, returning a JSON as a string, following the structure given
below. For the dates, consider just the year, month and day in the format YYYY-MM-DD.
Do not add any markup or markdown to the result.

results = {
  'domain_name'
  'expiration_date'
  'creation_date'
  'registrar'
  'name_servers'
  'registrant_contact'
  'admin_contact'
  'tech_contact'
  'status'
}
`

	// Prepare the request
	request := ChatCompletionRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: roleSystem,
			},
			{
				Role:    "user",
				Content: whoisText,
			},
		},
		MaxTokens:   intPtr(512),
		Temperature: float64Ptr(0.7),
		TopP:        float64Ptr(0.95),
	}

	// Send request
	response, err := c.ChatCompletion(request)
	if err != nil {
		return nil, err
	}

	var parsedWhois map[string]interface{}
	err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &parsedWhois)

	return parsedWhois, nil
}
