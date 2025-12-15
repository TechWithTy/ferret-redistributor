package main

import "testing"

func TestMessageLogID_IsDeterministic(t *testing.T) {
	a := messageLogID("Response", "bot123", "group456", "msg789")
	b := messageLogID("Response", "bot123", "group456", "msg789")
	if a != b {
		t.Fatalf("expected deterministic id, got %q vs %q", a, b)
	}
}

func TestMessageLogID_DiffersWhenInputsDiffer(t *testing.T) {
	a := messageLogID("Response", "bot123", "group456", "msg789")
	b := messageLogID("Outbound", "bot123", "group456", "msg789")
	if a == b {
		t.Fatalf("expected different ids when direction differs, got %q", a)
	}
}


