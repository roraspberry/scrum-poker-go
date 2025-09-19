package main

import (
	"encoding/json"
	"log"
	"net/http"
	"scrum-poker-go/internal/session"
	"strings"
	"sync"
)

// API holds in-memory sessions and a mutex.
type API struct {
	mu       sync.RWMutex
	sessions map[string]*session.Session
}

func NewAPI() *API {
	return &API{
		sessions: make(map[string]*session.Session),
	}
}

// Helper to write JSON response.
func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// POST /sessions to create session.
func (a *API) createSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Define a struct to hold the incoming request.
	var req struct {
		Title string `json:"title"`
	}

	// Decode JSON body into request body.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body.", http.StatusBadRequest)
		return
	}

	// Create a new session with the title.
	s := session.NewSession(req.Title)

	// Save session into API map.
	a.mu.Lock()
	a.sessions[s.ID] = s
	a.mu.Unlock()

	// Encode and send back JSON response.
	writeJSON(w, http.StatusCreated, s)
}

// Handles paths starting with /sessions/{id} or /sessions/{id}/join
func (a *API) sessionHandler(w http.ResponseWriter, r *http.Request) {
	// Path: /sessions/{id} or /sessions/{id}/join
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// parts[0] == "sessions", parts[1] = id, optional parts[2] = "join"
	if len(parts) < 2 || parts[0] != "sessions" {
		http.NotFound(w, r)
		return
	}

	id := parts[1]

	a.mu.RLock()
	s, ok := a.sessions[id]
	a.mu.RUnlock()

	if !ok {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	// /sessions/{id} - GET(view)
	if len(parts) == 2 && r.Method == http.MethodGet {
		// Convert players map into slice for JSON.
		players := make([]session.Player, 0, len(s.Players))
		for _, p := range s.Players {
			players = append(players, p)
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"id":      s.ID,
			"title":   s.Title,
			"players": players,
		})
		return
	}

	// /sessions/{id}/join - POST
	if len(parts) == 3 && parts[2] == "join" && r.Method == http.MethodPost {
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
			return
		}
		player, err := s.AddPlayer(req.Name)
		if err != nil {
			http.Error(w, "cannot add player: "+err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusCreated, player)
		return
	}

	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func main() {
	api := NewAPI()

	http.HandleFunc("/sessions", api.createSessionHandler)
	http.HandleFunc("/sessions/", api.sessionHandler)

	addr := ":8080"
	log.Printf("API server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
