package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/artemzi/olx-parser/version"
	"github.com/gorilla/mux"
)

func TestHealthHandler(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/healthz", healthz).Methods("GET")
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	if trw.Body.String() != "Ok" {
		t.Error("Expected", "Ok", "got", trw.Body.String())
	}
	if trw.Code != 200 {
		t.Error("Expected status:", 200, "got", trw.Body.String())
	}
}

func TestInfoHandler(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/info", info).Methods("GET")
	req, err := http.NewRequest("GET", "/info", nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	inf := "{\"commit\":\"" + version.COMMIT + "\",\"repo\":\"" + version.REPO +
		"\",\"version\":\"" + version.RELEASE + "\"}"
	if trw.Body.String() != inf {
		t.Error("Expected", inf, "got", trw.Body.String())
	}
	if trw.Code != 200 {
		t.Error("Expected status:", 200, "got", trw.Body.String())
	}
}

func TestRootHandler(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/", root).Methods("GET")
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	inf := fmt.Sprintf("PARSER v%s\n", version.RELEASE)
	if trw.Body.String() != inf {
		t.Error("Expected", inf, "got", trw.Body.String())
	}
	if trw.Code != 200 {
		t.Error("Expected status:", 200, "got", trw.Body.String())
	}
}
