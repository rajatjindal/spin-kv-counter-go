package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	spinhttp "github.com/fermyon/spin-go-sdk/http"
	kv "github.com/fermyon/spin-go-sdk/kv"
)

type Counter struct {
	Count int `json:"count"`
}

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		store, err := kv.OpenStore("default")
		if err != nil {
			http.Error(w, "failed to open store", http.StatusInternalServerError)
			return
		}

		counter, err := getJson[Counter](store, "counter")
		if err != nil {
			http.Error(w, "failed to get counter from kv", http.StatusInternalServerError)
			return
		}

		counter.Count += 1

		updatedValue, err := setJson(store, "counter", counter)
		if err != nil {
			http.Error(w, "failed to get counter from kv", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(updatedValue)
	})
}

func main() {}

func getJson[T any](store *kv.Store, key string) (T, error) {
	var val T
	exists, err := store.Exists("counter")
	if err != nil {
		return val, err
	}

	if !exists {
		return val, nil
	}

	value, err := store.Get(key)
	if err != nil {
		return val, fmt.Errorf("failed to get value for key %q from kv store", key)
	}

	err = json.Unmarshal(value, &val)
	if err != nil {
		return val, fmt.Errorf("failed to unmarshal into struct Counter")
	}

	return val, nil
}

func setJson[T any](store *kv.Store, key string, value T) ([]byte, error) {
	updatedValue, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated value into raw bytes")
	}

	err = store.Set(key, updatedValue)
	if err != nil {
		return nil, fmt.Errorf("failed to update value in kv")
	}

	return updatedValue, nil
}
