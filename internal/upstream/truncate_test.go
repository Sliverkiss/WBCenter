package upstream

import (
	"strings"
	"testing"
	"unicode/utf8"
)

// TestTruncateKeepsValidUTF8 截断必须停在字符边界。
//
// 上游错误文案大量是中文（「积分不足」「请求过于频繁」「今日已签到」…），而 truncate
// 的 n 是字节数：按字节切会把最后一个汉字劈成半个，产出**非法 UTF-8**。这段文案随后会
// 进入 Error.Msg / JSON 响应，非法字节在 json 编码时变成 U+FFFD，用户看到的就是乱码，
// 正好发生在最需要看清原因的错误路径上。
func TestTruncateKeepsValidUTF8(t *testing.T) {
	cases := []struct {
		name string
		in   string
		n    int
	}{
		{"中文-砍在字符中间", "签到失败：上游返回积分不足，请稍后重试", 20},
		{"中文-砍在第一个字后", "积分不足", 7},
		{"中文-长文案", "请求过于频繁，请稍后再试", 16},
		{"中英混合", "HTTP 401: token 已失效，请重新登录", 25},
		{"emoji", "旅行失败🚀请重试", 13},
		{"中文-边界正好", "积分不足", 6},
	}
	for _, c := range cases {
		got := truncate(c.in, c.n)
		if !utf8.ValidString(got) {
			t.Errorf("%s: truncate(%q, %d) 产出非法 UTF-8：%q（字节 % x）", c.name, c.in, c.n, got, []byte(got))
		}
		if len(got) > c.n {
			t.Errorf("%s: 截断后 %d 字节 > 上限 %d（%q）", c.name, len(got), c.n, got)
		}
		if !strings.HasPrefix(c.in, got) {
			t.Errorf("%s: 截断结果不是原串前缀：%q", c.name, got)
		}
	}
}

// TestTruncateStillFoldsAndKeepsShortText 原有的换行折叠与「不超限就不动」行为不能变。
func TestTruncateStillFoldsAndKeepsShortText(t *testing.T) {
	if got := truncate("积分不足\n请充值", 100); got != "积分不足 请充值" {
		t.Errorf("换行折叠失效：%q", got)
	}
	if got := truncate("短", 100); got != "短" {
		t.Errorf("未超限不应改动：%q", got)
	}
}
