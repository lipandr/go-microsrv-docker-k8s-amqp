package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, _ *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the Broker",
		Data:    nil,
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}
