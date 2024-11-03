package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		scores: map[string]int{
			"ジョンさん": 7,
			"vim":   1991,
		},
		winCalls: nil,
	}
	server := &PlayerServer{&store}

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
	server := PlayerServer{&store}

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

func assertStatus(want, got int, t *testing.T) {
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
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}
