package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL string
	Model   string
}

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type AnalysisResult struct {
	Summary     string `json:"summary"`
	RootCause   string `json:"root_cause"`
	Remediation string `json:"remediation"`
}

func NewClient(baseURL, model string) *Client {
	return &Client{
		BaseURL: baseURL,
		Model:   model,
	}
}

func (c *Client) Analyze(title, stackTrace, level string, count int) (*AnalysisResult, error) {
	prompt := fmt.Sprintf(`You are analyzing a production error. Provide a structured analysis.

	Error Details:
	- Title: %s
	- Level: %s
	- Occurrences: %d
	- Stack Trace:
	%s

	Provide your analysis in JSON format with these fields:
	{
		"summary": "2-3 sentence overview",
		"root_cause": "likely technical cause",
		"remediation": "recommended fix"
	}

	Respond ONLY with valid JSON, no markdown, no extra text.`, title, level, count, stackTrace)

	req := GenerateRequest{
		Model: c.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(req)
	if err != nil{
		return nil, err
	}

	resp, err := http.Post(
		c.BaseURL+"/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil{
		return nil, fmt.Errorf("ollama request failed: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}

	var genResp GenerateResponse
	err = json.Unmarshal(body, &genResp)
	if err != nil{
		return nil, err
	}

	var result AnalysisResult
	err = json.Unmarshal([]byte(genResp.Response), &result)
	if err != nil{
		result = AnalysisResult{
			Summary:     genResp.Response,
			RootCause:   "Unable to determine",
			Remediation: "Manual investigation required",
		}
	}

	return &result, nil
}