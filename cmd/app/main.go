package main

import (
	"fmt"
	"scrum-poker-go/internal/session"
)

func main() {
	s := session.NewSession("Sprint 1 Planning")
	p, _ := s.AddPlayer("Big Tuna")
	fmt.Printf("Created session %s with player %s\n", s.Title, p.Name)
}
