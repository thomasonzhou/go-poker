package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() []Player {
	return s.league
}

func TestGetPlayers(t *testing.T) {
	store := StubPlayerStore{
		scores: map[string]int{
			"ジョンさん": 7,
			"vim":   1991,
		},
		winCalls: nil,
	}
	server := NewPlayerServer(&store)

	t.Run("returns ジョンさん's score", func(t *testing.T) {
		request := newGetScoreRequest("ジョンさん")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(http.StatusOK, response.Code, t)
		assertResponseBody("7", response.Body.String(), t)
	})

	t.Run("gets Vim's score", func(t *testing.T) {
		request := newGetScoreRequest("vim")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(http.StatusOK, response.Code, t)
		assertResponseBody("1991", response.Body.String(), t)
	})

	t.Run("returns 404 on missing player", func(t *testing.T) {
		request := newGetScoreRequest("missing")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(http.StatusNotFound, response.Code, t)

	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		scores:   map[string]int{},
		winCalls: nil,
	}
	server := NewPlayerServer(&store)

	t.Run("stores wins on POST", func(t *testing.T) {
		player := "金太郎"
		response := httptest.NewRecorder()

		server.ServeHTTP(response, newPostWinRequest(player))

		assertStatus(http.StatusAccepted, response.Code, t)

		if len(store.winCalls) != 1 {
			t.Errorf("expected winCalls of length 1, but got length %d", len(store.winCalls))
		}

		if store.winCalls[0] != player {
			t.Errorf("expected player %q, got %q", player, store.winCalls[0])
		}
	})
}

func TestLeague(t *testing.T) {
	league := []Player{
		{"Genki", 2},
		{"Bulbasaur", 001},
		{"of Legends", 2009},
	}

	store := StubPlayerStore{league: league}
	server := NewPlayerServer(&store)

	response := httptest.NewRecorder()

	server.ServeHTTP(response, newLeagueRequest())

	got := getLeagueFromResponse(response, t)
	assertLeaguesMatch(league, got, t)
	assertContentType(response, jsonContentType, t)

	assertStatus(http.StatusOK, response.Code, t)
}

func assertLeaguesMatch(league []Player, got []Player, t *testing.T) {
	if !slices.Equal(league, got) {
		t.Errorf("Wanted %v got %v", league, got)
	}
}

func assertContentType(response *httptest.ResponseRecorder, want string, t *testing.T) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("Response content-type was not json, got %v", response.Result().Header)
	}
}

func newLeagueRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return request
}

func getLeagueFromResponse(response *httptest.ResponseRecorder, t *testing.T) (league []Player) {
	t.Helper()

	err := json.NewDecoder(response.Body).Decode(&league)
	if err != nil {
		t.Fatalf("Unable to parse %q into JSON, %v", response.Body, err)
	}
	return league
}

func assertStatus(want, got int, t *testing.T) {
	t.Helper()
	if got != want {
		t.Errorf("got status %d, want status %d", got, want)
	}
}

func newGetScoreRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/player/%s", player), nil)
	return request
}

func newPostWinRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/player/%s", player), nil)
	return request
}

func assertResponseBody(want, got string, t *testing.T) {
	t.Helper()
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}
