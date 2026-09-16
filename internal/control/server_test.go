package control

import "testing"

func TestRedactTaskDetailRemovesSecretsRecursively(t *testing.T) {
	got := redactTaskDetail(map[string]any{
		"name": "daily task", "accessToken": "must-not-leak",
		"nested": map[string]any{"cookie": "must-not-leak", "enabled": true},
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
}

// TestRedactTaskDetailRecursesIntoArrays 数组元素内部同样要脱敏。
//
// 上游任务列表就是数组套对象（{"tasks":[{...,"token":...}]}），旧实现
// 对 []any 只做 append(child[:100]...) 整体拷贝，数组元素里的敏感键
// 既不被删也不被检查，直接透出到管理视图——与函数声明的目的相悖。
func TestRedactTaskDetailRecursesIntoArrays(t *testing.T) {
	got := redactTaskDetail(map[string]any{
		"tasks": []any{
			map[string]any{"name": "t1", "accessToken": "must-not-leak", "cookie": "must-not-leak"},
			map[string]any{"name": "t2", "nested": map[string]any{"refreshToken": "must-not-leak"}},
		},
		// 数组套数组套对象：递归必须无深度限制。
		"deep": []any{[]any{map[string]any{"secret": "must-not-leak", "keep": "ok"}}},
	})

	tasks, ok := got["tasks"].([]any)
	if !ok || len(tasks) != 2 {
		t.Fatalf("tasks=%#v want 长度 2 的 []any", got["tasks"])
	}
	first, ok := tasks[0].(map[string]any)
	if !ok {
		t.Fatalf("tasks[0]=%#v want map[string]any", tasks[0])
	}
	if _, leaked := first["accessToken"]; leaked {
		t.Error("数组元素内的 accessToken 泄漏")
	}
	if _, leaked := first["cookie"]; leaked {
		t.Error("数组元素内的 cookie 泄漏")
	}
	if first["name"] != "t1" {
		t.Errorf("数组元素内的非敏感字段应保留，name=%v", first["name"])
	}
	second, ok := tasks[1].(map[string]any)
	if !ok {
		t.Fatalf("tasks[1]=%#v want map[string]any", tasks[1])
	}
	nested, ok := second["nested"].(map[string]any)
	if !ok {
		t.Fatalf("tasks[1].nested=%#v want map[string]any", second["nested"])
	}
	if _, leaked := nested["refreshToken"]; leaked {
		t.Error("数组元素内嵌套 map 的 refreshToken 泄漏")
	}

	deep, ok := got["deep"].([]any)
	if !ok || len(deep) != 1 {
		t.Fatalf("deep=%#v want 长度 1 的 []any", got["deep"])
	}
	inner, ok := deep[0].([]any)
	if !ok || len(inner) != 1 {
		t.Fatalf("deep[0]=%#v want 长度 1 的 []any", deep[0])
	}
	leaf, ok := inner[0].(map[string]any)
	if !ok {
		t.Fatalf("deep[0][0]=%#v want map[string]any", inner[0])
	}
	if _, leaked := leaf["secret"]; leaked {
		t.Error("数组套数组内对象的 secret 泄漏")
	}
	if leaf["keep"] != "ok" {
		t.Errorf("数组套数组内的非敏感字段应保留，keep=%v", leaf["keep"])
	}
}

// TestRedactTaskDetailArrayCap 数组截断上限在递归之后依然生效（保留旧行为）。
func TestRedactTaskDetailArrayCap(t *testing.T) {
	items := make([]any, 0, 150)
	for i := 0; i < 150; i++ {
		items = append(items, map[string]any{"name": "t", "token": "must-not-leak"})
	}
	got := redactTaskDetail(map[string]any{"tasks": items})
	tasks, ok := got["tasks"].([]any)
	if !ok {
		t.Fatalf("tasks=%#v want []any", got["tasks"])
	}
	if len(tasks) != maxRedactArrayItems {
		t.Errorf("数组应截断到 %d 项，实际 %d", maxRedactArrayItems, len(tasks))
	}
	if _, leaked := tasks[0].(map[string]any)["token"]; leaked {
		t.Error("截断后仍应脱敏，token 泄漏")
	}
}
