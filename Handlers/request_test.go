package main

import (
	"io/ioutil"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHPHandler(t *testing.T) {
	t.Run("it does something", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test?name=james", nil)
		w := httptest.NewRecorder()
		HPHandler(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expect error to be nil but got %v", err)
		}
		if string(data) != "James" {
			t.Errorf("expected James but got %v", string(data))
		}
	})
}