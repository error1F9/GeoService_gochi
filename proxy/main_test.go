package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReverseProxy(t *testing.T) {
	localBackendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Local Backend response")
	}))
	defer localBackendServer.Close()

	outBackendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Outsource Backend response")
	}))
	defer outBackendServer.Close()

	port := strings.Split(outBackendServer.URL, ":")[2]

	rp := NewReverseProxy("localhost", port)
	handler := rp.ReverseProxy(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, localBackendServer.URL, http.StatusTemporaryRedirect)
	}))

	tests := []struct {
		name         string
		path         string
		expected     string
		expectedCode int
	}{
		{
			name:         "Redirected request",
			path:         localBackendServer.URL + "/",
			expected:     "Local Backend response",
			expectedCode: 200,
		},
		{
			name:         "API request",
			path:         localBackendServer.URL + "/api/",
			expected:     "API response",
			expectedCode: 307,
		},
		{
			name:         "Main server request",
			path:         outBackendServer.URL + "/",
			expected:     "Backend response",
			expectedCode: 200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", test.path, nil)

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != test.expectedCode {
				t.Errorf("Expected %d, got %d", test.expectedCode, rec.Code)
			}

		})
	}
}
