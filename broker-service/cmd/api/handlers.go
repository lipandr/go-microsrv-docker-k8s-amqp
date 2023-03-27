package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, _ *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the Broker",
		Data:    nil,
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	if err := app.readJSON(w, r, &requestPayload); err != nil {
		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		_ = app.errorJSON(w, errors.New("unknown action"), http.StatusBadRequest)
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create json we'll send to the authentication endpoint
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest(
		http.MethodPost,
		"http://authentication-service/authenticate",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	// make sure the service returns a valid status code
	if response.StatusCode == http.StatusUnauthorized {
		_ = app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		_ = app.errorJSON(w, errors.New("error calling authentication service"))
		return
	}
	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}
	if jsonFromService.Error {
		_ = app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}
	// send the user back to the client
	payload := jsonResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    jsonFromService.Data,
	}

	_ = app.writeJSON(w, http.StatusAccepted, payload)
}
