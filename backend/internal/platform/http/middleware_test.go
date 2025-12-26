package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddleware_SingleOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CORSMiddleware("http://localhost:3000"))
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	router.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Errorf("Expected Access-Control-Allow-Origin to be 'http://localhost:3000', got '%s'", w.Header().Get("Access-Control-Allow-Origin"))
	}

	if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Error("Expected Access-Control-Allow-Credentials to be 'true'")
	}
}

func TestCORSMiddleware_MultipleOrigins_MatchingOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CORSMiddleware("http://localhost:3000,http://localhost:5173"))
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Test first origin
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	router.ServeHTTP(w, req)

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "http://localhost:3000" {
		t.Errorf("Expected Access-Control-Allow-Origin to be 'http://localhost:3000', got '%s'", origin)
	}

	// Test second origin
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	router.ServeHTTP(w, req)

	origin = w.Header().Get("Access-Control-Allow-Origin")
	if origin != "http://localhost:5173" {
		t.Errorf("Expected Access-Control-Allow-Origin to be 'http://localhost:5173', got '%s'", origin)
	}
}

func TestCORSMiddleware_MultipleOrigins_NoMatchingOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CORSMiddleware("http://localhost:3000,http://localhost:5173"))
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Request from non-allowed origin should fall back to first origin
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://evil.com")
	router.ServeHTTP(w, req)

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "http://localhost:3000" {
		t.Errorf("Expected fallback to first origin 'http://localhost:3000', got '%s'", origin)
	}
}

func TestCORSMiddleware_WildcardOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CORSMiddleware("*"))
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://any-origin.com")
	router.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Expected Access-Control-Allow-Origin to be '*', got '%s'", w.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestCORSMiddleware_OptionsPreflightRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CORSMiddleware("http://localhost:3000"))
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected OPTIONS request to return 204, got %d", w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("Expected Access-Control-Allow-Methods header to be set")
	}

	if w.Header().Get("Access-Control-Allow-Headers") == "" {
		t.Error("Expected Access-Control-Allow-Headers header to be set")
	}
}

func TestCORSMiddleware_OriginsWithSpaces(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	// Test that spaces around origins are trimmed
	router.Use(CORSMiddleware("http://localhost:3000 , http://localhost:5173"))
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	router.ServeHTTP(w, req)

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "http://localhost:5173" {
		t.Errorf("Expected trimmed origin 'http://localhost:5173', got '%s'", origin)
	}
}

func TestRequestLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestLoggerMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Logger middleware should not affect response
	if w.Body.String() != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", w.Body.String())
	}
}
