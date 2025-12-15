package notion

import (
	"testing"
	"time"
)

func TestNextRun(t *testing.T) {
	base := time.Date(2025, 12, 15, 10, 30, 0, 0, time.UTC)
	tests := []struct {
		rec      string
		wantOk   bool
		wantTime time.Time
	}{
		{"None", false, time.Time{}},
		{"", false, time.Time{}},
		{"Daily", true, base.Add(24 * time.Hour)},
		{"Every 7 days", true, base.Add(7 * 24 * time.Hour)},
		{"Every 2 weeks", true, base.Add(14 * 24 * time.Hour)},
		{"Monthly", true, base.AddDate(0, 1, 0)},
	}
	for _, tt := range tests {
		got, ok := NextRun(base, tt.rec)
		if ok != tt.wantOk {
			t.Fatalf("rec=%q ok=%v want %v", tt.rec, ok, tt.wantOk)
		}
		if tt.wantOk && !got.Equal(tt.wantTime) {
			t.Fatalf("rec=%q got=%s want=%s", tt.rec, got, tt.wantTime)
		}
	}
}


