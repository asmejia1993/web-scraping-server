package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func (hf handlerFranchises) respond(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	if data != nil {
		// Encode the actual response data, not respData
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			// Log the error if encoding fails
			http.Error(w, "Could not encode data in json", http.StatusInternalServerError)
			return
		}
	}
}

// it reads to the memory.
func (hf handlerFranchises) readRequestBody(r *http.Request) ([]byte, error) {
	// Read the content
	var bodyBytes []byte
	var err error
	if r.Body != nil {
		bodyBytes, err = io.ReadAll(r.Body)
		if err != nil {
			err := errors.New("could not read request body")
			return nil, err
		}
	}
	return bodyBytes, nil
}

func (hf handlerFranchises) decode(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (hf handlerFranchises) restoreRequestBody(r *http.Request, bodyBytes []byte) {
	// Restore the io.ReadCloser to its original state
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
}
