package main

import (
	"log-service/data"
	"net/http"
	"time"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json to var
	var payload JSONPayload
	_ = app.readJSON(w, r, &payload)

	// insert payload into db
	event := data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := app.Models.LogEntry.Insert(event); err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "Log entry created successfully!",
	}
	_ = app.writeJSON(w, http.StatusCreated, resp)
}
