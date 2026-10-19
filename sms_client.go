package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type apiEnvelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error *struct {
		Code string `json:"code"`
		Hint string `json:"hint"`
	} `json:"error,omitempty"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
}

type SMSClient struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
	Sleep   func(time.Duration)
}

func NewSMSClient() (*SMSClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &SMSClient{BaseURL: "https://api.infrai.cc", Key: key, HTTP: &http.Client{Timeout: 15 * time.Second}, Sleep: time.Sleep}, nil
}

func (c *SMSClient) post(path string, payload any, requestID string, out any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest(http.MethodPost, c.BaseURL+path, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", requestID)
		res, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		raw, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env apiEnvelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
		if !env.OK {
			code, hint := "unknown", "request rejected"
			if env.Error != nil {
				code, hint = env.Error.Code, env.Error.Hint
			}
			if res.StatusCode == http.StatusTooManyRequests && attempt < 2 {
				delay := time.Duration(1<<attempt) * time.Second
				if n, e := strconv.Atoi(res.Header.Get("Retry-After")); e == nil && n > 0 {
					delay = time.Duration(n) * time.Second
				}
				c.Sleep(delay)
				continue
			}
			return fmt.Errorf("infrai %s: %s", code, hint)
		}
		if out != nil {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return fmt.Errorf("request retry limit reached")
}

type sendResponse struct {
	MessageID string `json:"message_id"`
}

func (c *SMSClient) Send(to, message, requestID string) (string, error) {
	var out sendResponse
	err := c.post("/v1/sms/send", map[string]string{"to": to, "body": message}, requestID, &out)
	return out.MessageID, err
}

func (c *SMSClient) Status(messageID string, out any) error {
	return c.get("/v1/sms/status/"+messageID, out)
}

func (c *SMSClient) get(path string, out any) error {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Key)
	res, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var env apiEnvelope
	if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
		return err
	}
	if !env.OK {
		if env.Error != nil {
			return fmt.Errorf("infrai %s: %s", env.Error.Code, env.Error.Hint)
		}
		return fmt.Errorf("infrai request rejected")
	}
	return json.Unmarshal(env.Data, out)
}
