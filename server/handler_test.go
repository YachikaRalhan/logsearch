package server

import (
	"bytes"
	"encoding/json"
	"logpuller/pkg/model"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestRetrieveLogFiles(t *testing.T) {
	// Test case to ensure log files are retrieved correctly within a specified time range
	from := time.Date(2024, time.May, 4, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, time.May, 4, 1, 0, 0, 0, time.UTC)

	logFiles, err := retrieveLogFiles(from, to)
	os.Setenv("AWS_S3_BUCKET_NAME", "<add-bucket-name>")
	os.Setenv("AWS_REGION", "<add-region-name>")

	if err != nil {
		t.Errorf("RetrieveLogFiles returned error: %v", err)
	}

	if len(logFiles) < 1 {
		t.Errorf("Expected log files, got %d", len(logFiles))
	}
}

func TestLogSearchHandler(t *testing.T) {
	// Test case for the logsearch handler

	// Create a mock request body
	from := time.Date(2024, time.May, 4, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, time.May, 4, 1, 0, 0, 0, time.UTC)
	requestData := model.LogSearchRequest{
		From:          from,
		To:            to,
		SearchKeyword: "Server restarted.",
	}
	requestBody, _ := json.Marshal(requestData)
	os.Setenv("AWS_S3_BUCKET_NAME", "<add-bucket-name>")
	os.Setenv("AWS_REGION", "<add-region-name>")
	// Create a mock HTTP request
	req, err := http.NewRequest("POST", "/logsearch", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the logsearch handler
	handler := http.HandlerFunc(logsearch)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"logLines":["2024-05-04 01:50:21 - INFO - Server restarted."]}`
	res := strings.TrimSpace(rr.Body.String())
	if res != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
