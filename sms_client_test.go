package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendRequestFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/sms/send" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var fields map[string]any
		if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
			t.Errorf("decode request: %v", err)
		}
		if len(fields) != 2 || fields["to"] != "test-number" || fields["body"] != "test reminder" {
			t.Errorf("unexpected send fields: %v", fields)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"data":{"message_id":"test-id"}}`))
	}))
	defer server.Close()

	client := &SMSClient{BaseURL: server.URL, Key: "test-key", HTTP: server.Client()}
	id, err := client.Send("test-number", "test reminder", "test-request")
	if err != nil || id != "test-id" {
		t.Fatalf("Send() = %q, %v", id, err)
	}
}
