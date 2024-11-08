package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestRecordingWinsAndRetrievingWins(t *testing.T) {

	store := NewInMemoryStore()
	server := NewPlayerServer(store)
	const player = "Mikey"
	const writeCount = 3

	for range writeCount {
		server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	}

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))

		assertStatus(http.StatusOK, response.Code, t)
		assertResponseBody(strconv.Itoa(writeCount), response.Body.String(), t)
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())

		assertStatus(http.StatusOK, response.Code, t)
		assertContentType(response, jsonContentType, t)

		league := getLeagueFromResponse(response, t)
		assertLeaguesMatch(league, []Player{{player, 3}}, t)
	})
}
