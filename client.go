package regulayer

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Client is the Regulayer SDK client
type Client struct {
	APIKey     string
	Endpoint   string
	HTTPClient *http.Client
}

// Config for initializing the client
type Config struct {
	APIKey   string
	Endpoint string
	Demo     bool
}

// Decision payload
type Decision struct {
	DecisionID string
	System     string
	RiskLevel  string
	ModelName  string
	Input      map[string]interface{}
	Output     map[string]interface{}
	Metadata   map[string]interface{}
}

// NewClient creates a new Regulayer client
func NewClient(config Config) (*Client, error) {
	if config.APIKey == "" {
		return nil, errors.New("APIKey is required")
	}

	endpoint := config.Endpoint
	if endpoint == "" {
		endpoint = "https://api.regulayer.tech"
	}

	return &Client{
		APIKey:     config.APIKey,
		Endpoint:   endpoint,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// generateUUIDv4 creates a secure random UUIDv4
func generateUUIDv4() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// hashJSON computes the SHA-256 hash of a JSON encoded map
func hashJSON(data interface{}) string {
	if data == nil {
		return ""
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:])
}

// RecordDecision sends a decision to the Regulayer gateway
func (c *Client) RecordDecision(d Decision) error {
	decisionID := d.DecisionID
	if decisionID == "" {
		decisionID = generateUUIDv4()
	}

	riskLevel := d.RiskLevel
	if riskLevel == "" {
		riskLevel = "standard"
	}

	modelName := d.ModelName
	if modelName == "" {
		modelName = "default"
	}

	now := time.Now().UTC()

	payload := map[string]interface{}{
		"event_version":   "2.0",
		"event_state":     "completed",
		"decision_id":     decisionID,
		"system_name":     d.System,
		"risk_level":      riskLevel,
		"model_name":      modelName,
		"input_hash":      hashJSON(d.Input),
		"output_hash":     hashJSON(d.Output),
		"input":           d.Input,
		"output":          d.Output,
		"metadata":        d.Metadata,
		"start_timestamp": now.Format(time.RFC3339),
		"end_timestamp":   now.Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.Endpoint+"/v1/decisions", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("X-Regulayer-Api-Key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", decisionID)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: status code %d", resp.StatusCode)
	}

	return nil
}
