package session

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Player represents a participant in a session.
type Player struct {
	ID   string
	Name string
}

// Session holds the state of a poker session.
type Session struct {
	ID        string
	Title     string
	CreatedAt time.Time
	Players   map[string]Player
	mu        sync.Mutex
}

// NewSession creates a new planning poker session.
func NewSession(title string) *Session {
	return &Session{
		ID:        uuid.New().String(),
		Title:     title,
		CreatedAt: time.Now(),
		Players:   make(map[string]Player),
	}
}

// AddPlayer adds a player to the session.
func (s *Session) AddPlayer(name string) (Player, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Prevent duplicate names (optional rule)
	for _, p := range s.Players {
		if p.Name == name {
			return Player{}, errors.New("player with this name already exists")
		}
	}

	player := Player{
		ID:   uuid.New().String(),
		Name: name,
	}
	s.Players[player.ID] = player
	return player, nil
}
