package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSuccessResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := map[string]string{"message": "success"}
	SuccessResponse(c, testData)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'data' field in response")
	}

	if data["message"] != "success" {
		t.Errorf("Expected message 'success', got '%v'", data["message"])
	}
}

func TestListResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testData := []string{"item1", "item2"}
	meta := map[string]interface{}{
		"total": 2,
		"page":  1,
	}
	ListResponse(c, testData, meta)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if _, ok := response["data"]; !ok {
		t.Error("Expected 'data' field in response")
	}

	responseMeta, ok := response["meta"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'meta' field in response")
	}

	if responseMeta["total"].(float64) != 2 {
		t.Errorf("Expected total 2, got %v", responseMeta["total"])
	}
}

func TestErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid input", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	errorData, ok := response["error"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected 'error' field in response")
	}

	if errorData["code"] != "VALIDATION_ERROR" {
		t.Errorf("Expected error code 'VALIDATION_ERROR', got '%v'", errorData["code"])
	}

	if errorData["message"] != "Invalid input" {
		t.Errorf("Expected error message 'Invalid input', got '%v'", errorData["message"])
	}
}
