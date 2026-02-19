package atlas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

type Option func(*Client)

type Client struct {
	apiKey  string
	baseURL string
	enabled bool
	client  *http.Client
}

type Event struct {
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	StackTrace string    `json:"stack_trace"`
	Timestamp  time.Time `json:"timestamp"`
}

func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

func WithEnabled(enabled bool) Option {
	return func(c *Client) {
		c.enabled = enabled
	}
}

func NewClient(apiKey string, options ...Option) *Client{
	client := &Client{
		apiKey:  apiKey,
		baseURL: "http://localhost:8081",
		enabled: true,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
	for _, opt := range options{
		opt(client)
	}

	return client
}

func (c *Client) CaptureError(err error) {
	if !c.enabled || err == nil {
		return
	}

	event := Event{
		Level:      "error",
		Message:    err.Error(),
		StackTrace: string(debug.Stack()),
		Timestamp:  time.Now().UTC(),
	}

	c.send(event)
}

func (c *Client) CaptureMessage(message, level string) {
	if !c.enabled {
		return
	}

	event := Event{
		Level:      level,
		Message:    message,
		StackTrace: "",
		Timestamp:  time.Now().UTC(),
	}

	c.send(event)
}

func (c *Client) send(event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("Atlas: Failed to marshal event: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/ingest/events", bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Atlas: Failed to create request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Printf("Atlas: Network error: %v\n", err)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 202 {
		fmt.Printf("Atlas: Failed to send event (%d)\n", resp.StatusCode)
	}
}

func (c *Client) GinMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			err := recover(); 
			if err != nil {
				c.CaptureError(fmt.Errorf("panic: %v", err))
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		ctx.Next()
	}
}