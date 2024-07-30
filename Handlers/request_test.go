package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"go.mongodb.org/mongo-driver/bson"
)

var mockStore = &MockCharacterStore{
	Data: map[string]bson.M{
		"123": {"name": "Johnny Silverhand", "hp": 100},
	},
}

func TestHPHandler(t *testing.T) {
	t.Run("GET request - successful", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/HP?id=123", nil)
		rec := httptest.NewRecorder()
		HPHandler := makeHPHandler(mockStore)
		HPHandler(rec, req)

		res := w.Result()
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("expected no error but got %v", err)
		}

		expected := `{"name":"Johnny Silverhand","hp":100}`
		if string(data) != expected {
			t.Errorf("expected %s but got %s", expected, string(data))
		}
	})

	t.Run("GET request - name not found", func(t *testing.T) {
	})

	t.Run("POST request - successful", func(t *testing.T) {
	})

	t.Run("POST request - invalid JSON", func(t *testing.T) {
	})

	t.Run("POST request - name field missing", func(t *testing.T) {
	})

	t.Run("POST request - name not found", func(t *testing.T) {
	})
}
