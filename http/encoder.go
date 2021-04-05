package http

import (
	"encoding/json"
	"net/http"
)

type encoder struct {
	tracing bool
}

func newEncoder(tracing bool) *encoder {
	return &encoder{
		tracing: tracing,
	}
}

// errorResponse will encapsulate errors to be transferred over http.
type errorResponse struct {
	Message string      `json:"message"`
	Trace   interface{} `json:"trace,omitempty"`
	Trail   interface{} `json:"trail,omitempty"`
}

func (e *encoder) Response(w http.ResponseWriter, response []byte) error {
	return e.StatusResponse(w, response, http.StatusOK)
}

func (e *encoder) StatusResponse(w http.ResponseWriter, response []byte, status int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(response)
	return nil
}

func (e *encoder) Error(w http.ResponseWriter, err error, status int) {
	resp := errorResponse{
		Message: err.Error(),
	}
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
