package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestRecordingWinsAndRetrievingWins(t *testing.T) {
	store := NewInMemoryStore()
	server := PlayerServer{store}
	player := "Mikey"
	writeCount := 3

	for range writeCount {
		server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	}

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))

	assertStatus(http.StatusOK, response.Code, t)
	assertResponseBody(strconv.Itoa(writeCount), response.Body.String(), t)

}
