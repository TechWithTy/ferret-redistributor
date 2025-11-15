package workers

import "testing"

func TestTriggerMatcher(t *testing.T) {
    m := TriggerMatcher{Triggers: []string{"ready", "template"}}
    if _, ok := m.Match("I am READY!"); !ok { t.Fatal("should match") }
    if _, ok := m.Match("nope"); ok { t.Fatal("should not match") }
}
