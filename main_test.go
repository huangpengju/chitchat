package main

import (
	"chitchat/routes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFileServer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", routes.Index)

	w := httptest.NewRecorder()

	r, _ := http.NewRequest("GET", "/", nil)

	mux.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("Response code is%v", w.Code)
	}
	body := w.Body.String()
	if strings.Contains(body, "主页") == false {
		t.Errorf("Body does not contain 主页")
	}
}
