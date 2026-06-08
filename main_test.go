package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandlers is a table-driven test covering both routes. Table-driven tests
// are the Go idiom: one test function, many cases, easy to extend.
func TestHandlers(t *testing.T) {
	srv := newServer()

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string // exact body for plain responses; "" means skip check
	}{
		{name: "health is ok", path: "/healthz", wantStatus: http.StatusOK, wantBody: "ok"},
		{name: "root returns metadata", path: "/", wantStatus: http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			if tc.wantBody != "" && rec.Body.String() != tc.wantBody {
				t.Fatalf("body = %q, want %q", rec.Body.String(), tc.wantBody)
			}
		})
	}
}

// TestRootIsValidJSON ensures the root handler emits the expected JSON shape.
func TestRootIsValidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	newServer().ServeHTTP(rec, req)

	var got buildInfo
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if got.Service != "ga_practice" {
		t.Fatalf("service = %q, want %q", got.Service, "ga_practice")
	}
}
