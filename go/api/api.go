package api

import (
	"encoding/json"
	"net/http"
)

type Result struct {
	Data  any
	Error error
}

func (r *Result) MarshalJson() ([]byte, error) {
	if r == nil {
		return []byte("null"), nil
	}
	rj := struct {
		Data  any     `json:"data"`
		Error *string `json:"error"`
	}{
		Data: r.Data,
	}

	if r.Error != nil {
		s := r.Error.Error()
		rj.Error = &s
	}

	return json.Marshal(rj)
}

func writeResponseWithJson(w http.ResponseWriter, result Result) {
	w.Header().Set("Content-Type", "application/json")
	status := http.StatusOK
	var body []byte

	if result.Error != nil {
		status = http.StatusInternalServerError
	}
	if json, err := result.MarshalJson(); err != nil {
		body = []byte(err.Error())
		status = http.StatusInternalServerError
	} else {
		body = json
	}

	w.WriteHeader(status)
	w.Write([]byte(body))
}

func writeResponse(w http.ResponseWriter, result Result) {
	status := http.StatusOK
	var body []byte

	if result.Error != nil {
		body = []byte(result.Error.Error())
		status = http.StatusInternalServerError
	} else {
		body = []byte("Success")
	}

	w.WriteHeader(status)
	w.Write([]byte(body))
}
