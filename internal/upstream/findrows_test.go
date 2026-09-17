package upstream

import (
	"encoding/json"
	"fmt"
	"testing"
)

// TestFindObjectRowsIsDeterministic 同一份响应必须每次解析出同一批行。
//
// findObjectRows 的兜底分支是 `for _, child := range x`——遍历 map，而 Go 的 map 迭代
// 顺序是每次随机的（自 Go 1.12 起刻意随机化）。当响应里没有首选键、但同一层存在多个
// 数组时，挑中哪个数组就是随机的：面板的列表会在刷新之间抖动，而且因为 service 侧用
// UID:taskCode 作为「已见任务」的键，抖动还会让「新任务」通知反复触发。
func TestFindObjectRowsIsDeterministic(t *testing.T) {
	raw := []byte(`{"envelope":{"zzz":[{"code":"ZZZ"}],"aaa":[{"code":"AAA"}]}}`)
	var decoded any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}

	var first string
	for i := 0; i < 300; i++ {
		rows := findObjectRows(decoded)
		if len(rows) == 0 {
			t.Fatalf("第 %d 次没有解析出任何行", i)
		}
		got := fmt.Sprint(rows)
		if i == 0 {
			first = got
			continue
		}
		if got != first {
			t.Fatalf("同一份响应解析出不同结果（第 %d 次）：\n  got  %s\n  want %s\n"+
				"（兜底分支 range map 顺序随机，面板列表与「新任务」通知会随刷新抖动）", i, got, first)
		}
	}

	// 确定性口径：首选键清单之后，按键名字典序遍历。
	code, _ := findObjectRows(decoded)[0]["code"].(string)
	if code != "AAA" {
		t.Errorf("兜底选中 %q，期望字典序最小的 aaa（口径需可预期）", code)
	}
}

// TestFindObjectRowsPrefersStructuredKeys 首选键清单仍须优先于同级零星数组，
// 不因兜底确定性化而改变优先级。
func TestFindObjectRowsPrefersStructuredKeys(t *testing.T) {
	for i := 0; i < 50; i++ {
		raw := []byte(`{"aaa":[{"code":"AAA"}],"tasks":[{"code":"T"}],"zzz":[{"code":"ZZZ"}]}`)
		var decoded any
		if err := json.Unmarshal(raw, &decoded); err != nil {
			t.Fatal(err)
		}
		rows := findObjectRows(decoded)
		if len(rows) != 1 || rows[0]["code"] != "T" {
			t.Fatalf("第 %d 次未优先选 tasks: %v", i, rows)
		}
	}
}

// TestFindObjectRowsDescendsIntoNestedPreferredKey 嵌套形状（data/list 包裹）仍须解析出来。
func TestFindObjectRowsDescendsIntoNestedPreferredKey(t *testing.T) {
	raw := []byte(`{"code":0,"data":{"list":[{"id":"a"},{"id":"b"}]}}`)
	var decoded any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 50; i++ {
		rows := findObjectRows(decoded)
		if len(rows) != 2 || rows[0]["id"] != "a" || rows[1]["id"] != "b" {
			t.Fatalf("第 %d 次嵌套解析结果不对: %v", i, rows)
		}
	}
}
