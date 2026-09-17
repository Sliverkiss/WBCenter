package authstore

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// nestedDoc 造一份嵌套形凭证（与网关写出的形态一致）。
func nestedDoc(uid, nickname, token string) string {
	return `{"auth":{"accessToken":"` + token + `","refreshToken":"rt","expiresAt":1893456000,"domain":""},` +
		`"account":{"uid":"` + uid + `","nickname":"` + nickname + `"}}`
}

// TestGetFindsAccountInNonCanonicalFile 按 uid 取账号必须覆盖网关能加载的全部文件。
//
// 网关侧的加载口径是宽模式（AuthFileGlob = "workbuddy*.json"）：
// 不带连字符的文件（如 workbuddy_new.json）在网关侧是**正常账号**，会被加载并进池。
// 面板的 List 也用宽模式，因此这种账号一定会出现在界面上；但操作路径
// （签到 / 旅行 / 刷新凭据 / 积分详情）全部经由 Get(uid) 定位文件——
// 如果 Get 只认 workbuddy-<uid>.json，界面就会出现一个「列得出来、点什么都是
// 找不到文件」的账号，正是网关注释里说的那种两端口径对不上。
func TestGetFindsAccountInNonCanonicalFile(t *testing.T) {
	dir := t.TempDir()
	st, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "workbuddy_new.json"),
		[]byte(nestedDoc("uid-noncanon-1", "非规范文件名账号", "at-new")), 0o600); err != nil {
		t.Fatal(err)
	}

	rows, warnings := st.List()
	if len(warnings) != 0 {
		t.Fatalf("List 报警告: %v", warnings)
	}
	if len(rows) != 1 || rows[0].UID != "uid-noncanon-1" {
		t.Fatalf("List 未按宽模式列出该账号: %#v", rows)
	}

	got, err := st.Get("uid-noncanon-1")
	if err != nil {
		t.Fatalf("List 列出了 uid-noncanon-1，Get 却取不到: %v\n"+
			"（面板会显示一个所有操作都失败的账号）", err)
	}
	if got.AccessToken != "at-new" || got.Nickname != "非规范文件名账号" {
		t.Errorf("取到的账号内容不对: %#v", got)
	}
	if filepath.Base(got.FilePath) != "workbuddy_new.json" {
		t.Errorf("FilePath = %q，期望指回真实文件 workbuddy_new.json（刷新凭据要写回原文件）", got.FilePath)
	}
}

// TestGetPrefersCanonicalFile 同一个 uid 同时出现在规范文件与其它 workbuddy*.json 时，
// 规范文件优先——规则确定，不随目录顺序变化。
func TestGetPrefersCanonicalFile(t *testing.T) {
	dir := t.TempDir()
	st, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "workbuddy-uid-dup-1.json"),
		[]byte(nestedDoc("uid-dup-1", "规范文件", "at-canonical")), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "workbuddy_dup.json"),
		[]byte(nestedDoc("uid-dup-1", "另一个文件", "at-other")), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := st.Get("uid-dup-1")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.AccessToken != "at-canonical" {
		t.Errorf("应优先规范文件，实际取到 %q（FilePath=%q）", got.AccessToken, got.FilePath)
	}
	if base := filepath.Base(got.FilePath); base != "workbuddy-uid-dup-1.json" {
		t.Errorf("FilePath = %q，期望 workbuddy-uid-dup-1.json", base)
	}
}

// TestGetUnknownUIDStillErrors 回退逻辑不能「随便返回一个账号」。
func TestGetUnknownUIDStillErrors(t *testing.T) {
	dir := t.TempDir()
	st, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "workbuddy-someone.json"),
		[]byte(nestedDoc("uid-someone-1", "别的账号", "at-someone")), 0o600); err != nil {
		t.Fatal(err)
	}

	if got, err := st.Get("uid-missing-2"); err == nil {
		t.Fatalf("取不存在的 uid 竟然成功了: %#v", got)
	}
}

// TestGetStillRejectsInvalidUID 回退逻辑不能绕过 uid 校验（防目录穿越）。
func TestGetStillRejectsInvalidUID(t *testing.T) {
	dir := t.TempDir()
	st, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{"../../etc/passwd", "a/../../b", "", "x..y-padding"} {
		got, err := st.Get(bad)
		if err == nil {
			t.Errorf("非法 uid %q 竟然取到了: %#v", bad, got)
			continue
		}
		if !strings.Contains(err.Error(), "非法 uid") && !strings.Contains(err.Error(), "不存在") {
			t.Errorf("非法 uid %q 的错误信息不明确: %v", bad, err)
		}
	}
}
