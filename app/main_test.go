package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if resp["status"] != "ok" {
		t.Fatalf("expected status ok, got %s", resp["status"])
	}
}

func TestVersionHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	versionHandler(w, req)

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if resp["version"] != VERSION {
		t.Fatalf("expected version %s, got %s", VERSION, resp["version"])
	}
}

func TestEnvHandler(t *testing.T) {
	if err := os.Setenv("ENVIRONMENT", "test-env"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/env", nil)
	w := httptest.NewRecorder()

	envHandler(w, req)

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if resp["environment"] != "test-env" {
		t.Fatalf("expected test-env, got %s", resp["environment"])
	}
}

func TestPostConfigHandler(t *testing.T) {
	body := `{"name":"database_url","value":"postgres://example"}`
	req := httptest.NewRequest(http.MethodPost, "/config", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	postConfigHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if resp["name"] != "database_url" || resp["value"] != "postgres://example" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestGetConfigHandler(t *testing.T) {
	configStore.Lock()
	configStore.data["api_key"] = "12345"
	configStore.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/config/api_key", nil)
	w := httptest.NewRecorder()

	getConfigHandler(w, req)

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if resp["name"] != "api_key" || resp["value"] != "12345" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestDeleteConfigHandler(t *testing.T) {
	configStore.Lock()
	configStore.data["remove_me"] = "bye"
	configStore.Unlock()

	req := httptest.NewRequest(http.MethodDelete, "/config/remove_me", nil)
	w := httptest.NewRecorder()

	deleteConfigHandler(w, req)

	var resp map[string]bool
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if !resp["deleted"] {
		t.Fatalf("expected deleted=true, got %+v", resp)
	}

	configStore.RLock()
	_, exists := configStore.data["remove_me"]
	configStore.RUnlock()

	if exists {
		t.Fatalf("config still exists after delete")
	}
}
