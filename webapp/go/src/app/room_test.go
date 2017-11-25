package main

import "testing"

func TestGetRoomServer(t *testing.T) {
	s := getRoomServer("ho")
	expected := "app0122.isu7f.k0y.org"
	if s != expected {
		t.Errorf("out: %s, expected: %s", s, expected)
	}
}
