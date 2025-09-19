package session

import "testing"

func TestNewSession(t *testing.T) {
	s := NewSession("Sprint Planning")
	if s.Title != "Sprint Planning" {
		t.Errorf("expected title 'Sprint Planning', got '%s'", s.Title)
	}
	if s.ID == "" {
		t.Error("expected session ID to be generated")
	}
}

func TestAddPlayer(t *testing.T) {
	s := NewSession("Test Session")
	p, err := s.AddPlayer("Big Tuna")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "Big Tuna" {
		t.Errorf("expected player name 'Big Tuna', got '%s'", p.Name)
	}

	_, err = s.AddPlayer("Big Tuna")
	if err == nil {
		t.Error("expected error when adding duplicate player name")
	}
}
