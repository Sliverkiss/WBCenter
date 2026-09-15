package control

import (
	"errors"
	"testing"

	"workbuddy-control-center/internal/upstream"
)

func TestActivityReadErrorGivesActionableAccountStatus(t *testing.T) {
	cases := []struct {
		err  error
		want string
	}{
		{&upstream.Error{Kind: upstream.ErrSessionDead, Status: 401}, "活动读取失败，请刷新凭据或重新登录"},
		{&upstream.Error{Kind: upstream.ErrSoftRate, Status: 429}, "活动读取受限，请稍后再试"},
		{&upstream.Error{Kind: upstream.ErrNotFound, Status: 404}, "此账号暂不支持成长活动接口"},
		{errors.New("network failed"), "上游活动任务查询失败"},
	}
	for _, tc := range cases {
		if got := activityReadError(tc.err); got != tc.want {
			t.Fatalf("activityReadError(%v) = %q, want %q", tc.err, got, tc.want)
		}
	}
}

func TestReadableRunnerLineUsesChineseTaskSummary(t *testing.T) {
	names := map[string]string{"chat_5": "桌面端对话 5 次"}
	cases := []struct {
		line     string
		want     string
		terminal bool
	}{
		{"[task_runner] abc123 chat_5: accept 200 ok", "任务「桌面端对话 5 次」：已接受任务", false},
		{"[task_runner] abc123 chat_5: report 2/5 200 code=0 id=x", "任务「桌面端对话 5 次」：进度上报 2/5 成功", false},
		{"[task_runner] abc123 chat_5: claim 200 ok(credit=+100 energy=+5)", "任务「桌面端对话 5 次」：奖励领取成功", true},
	}
	for _, tc := range cases {
		got, code, terminal := readableRunnerLine(tc.line, names)
		if got != tc.want || code != "chat_5" || terminal != tc.terminal {
			t.Fatalf("readableRunnerLine(%q) = %q, %q, %v", tc.line, got, code, terminal)
		}
	}
}

func TestFormatDurationChinese(t *testing.T) {
	if got := formatDuration(125); got != "2 分 5 秒" {
		t.Fatalf("formatDuration(125) = %q", got)
	}
}
