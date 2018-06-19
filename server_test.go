package main

import (
	"fmt"
	"io"
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

func TestNotFoundHandler(t *testing.T) {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(logging(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, `Not Found`)
	}))

	req, err := http.NewRequest("GET",
		"/b6c8528975c1af6e14dac4e61",
		nil)
	if err != nil {
		t.Error(err)
	}
	trw := httptest.NewRecorder()
	r.ServeHTTP(trw, req)
	msg := "Not Found"
	if trw.Body.String() != msg {
		t.Error("Expected", msg, "got", trw.Body.String())
	}
	if trw.Code != 404 {
		t.Error("Expected status:", 404, "got", trw.Body.String())
	}
}
