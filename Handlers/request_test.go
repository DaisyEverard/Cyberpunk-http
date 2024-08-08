package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"go.mongodb.org/mongo-driver/bson"
	"bytes"
	"strings"
)

func getMockStore() *MockCharacterStore {
	return &MockCharacterStore{
		Data: []bson.M{
			{"_id": 123,"name": "Johnny Silverhand", "HP": 100,},
			{"_id": 13,"name": "Name", "HP": 20,},
		},
	}
}

func TestHPHandler(t *testing.T) {
	t.Run("GET request - successful", func(t *testing.T) {
		mockStore := getMockStore()
		req := httptest.NewRequest(http.MethodGet, "/HP?id=123", nil)
		rec := httptest.NewRecorder()
		HPHandler := makeDocumentHandler(mockStore)
		HPHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("expected no error but got %v", err)
		}

		jsonString := string(data)

		expectedPairs := []string{"\"name\":\"Johnny Silverhand\"", "\"HP\":100"}
		for _, pair := range expectedPairs {
			if !(strings.Contains(jsonString, pair)) {
				t.Errorf("didn't contain expected string '%s' in result %s",pair, string(data))
			}
		}
	})

	t.Run("GET request - id not found", func(t *testing.T) {
		mockStore := getMockStore()
		req := httptest.NewRequest(http.MethodGet, "/HP", nil)
		rec := httptest.NewRecorder()
		HPHandler := makeDocumentHandler(mockStore)
		HPHandler(rec, req)

		res := rec.Result()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected StatusBadRequest error but got %v", res.StatusCode)
		}
	})

	t.Run("POST request - successful", func(t *testing.T) {
		mockStore := getMockStore()
		jsonBody := []byte(`{"HP":5}`)
		bodyReader := bytes.NewReader(jsonBody)
		req := httptest.NewRequest(http.MethodPost, "/HP?id=123", bodyReader)

		rec := httptest.NewRecorder()
		HPHandler := makeDocumentHandler(mockStore)
		HPHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		_, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("expected no error but got %v", err)
		}
		res.Body.Close()

		expectedHP := 5
		actualHP := mockStore.Data[0]["HP"].(int)

		if actualHP != expectedHP {
			t.Errorf("expected %v but got %v.", expectedHP, actualHP)
		}
	})

	t.Run("POST request - invalid JSON", func(t *testing.T) {
	})

	t.Run("POST request - id field missing", func(t *testing.T) {
	})

	t.Run("POST request - id not found", func(t *testing.T) {
	})
}
