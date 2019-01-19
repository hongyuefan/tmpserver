package backserver

import (
	"testing"
)

func TestGetExchange(t *testing.T) {

	s := NewServer("a", 10, 10)

	t.Log(s.getExchange())
}
