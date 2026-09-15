package control

import "testing"

func TestRedactTaskDetailRemovesSecretsRecursively(t *testing.T) {
	got := redactTaskDetail(map[string]any{
		"name": "daily task", "accessToken": "must-not-leak",
		"nested": map[string]any{"cookie": "must-not-leak", "enabled": true},
		"steps":  []any{map[string]any{"secret": "must-not-leak", "title": "safe"}},
	})
	if _, ok := got["accessToken"]; ok {
		t.Fatal("token field leaked")
	}
	nested, ok := got["nested"].(map[string]any)
	if !ok || nested["enabled"] != true {
		t.Fatal("safe nested data missing")
	}
	if _, ok := nested["cookie"]; ok {
		t.Fatal("nested cookie field leaked")
	}
	steps, ok := got["steps"].([]any)
	if !ok || len(steps) != 1 {
		t.Fatal("safe array data missing")
	}
	step, ok := steps[0].(map[string]any)
	if !ok || step["title"] != "safe" {
		t.Fatal("safe nested array object missing")
	}
	if _, ok := step["secret"]; ok {
		t.Fatal("nested array secret leaked")
	}
}
