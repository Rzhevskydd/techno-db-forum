package utils

import (
	"encoding/json"
	"net/http"
)

type Dict map[string]string

func NewError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(Dict{"message": msg})
}

func CloseBody(w http.ResponseWriter, r *http.Request) {
	if err := r.Body.Close(); err != nil {
		NewError(w, http.StatusInternalServerError, err.Error())
	}
}

func ResponseJson(w http.ResponseWriter, code int, v interface{}) {
	w.WriteHeader(code)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}
