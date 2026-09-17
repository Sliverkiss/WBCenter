package fsutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestWriteFileAtomicSurvivesForeignTempWriter 面板与网关共用 auths 目录，而网关侧
// SaveAtomic 用的临时文件名就是 <path>.tmp（固定名）。面板若也用固定名，两个进程会写
// 同一个临时文件、再各自 rename，把彼此写了一半的字节混进正式凭证。
//
// 这里不靠时序去撞，而是**构造**一次交错：先替网关打开 path+".tmp" 写半截并保持句柄，
// 再调面板的 WriteFileAtomic，最后让网关把剩下半截写完。只要面板复用了同一个临时 inode，
// rename 之后网关的句柄仍指向正式文件，它的后续字节就会写进凭证里。
func TestWriteFileAtomicSurvivesForeignTempWriter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "workbuddy-uid-1.json")
	panelDoc := `{"auth":{"accessToken":"panel-token"},"account":{"uid":"uid-1"}}`

	// 网关同刻在写同一账号：已打开自己的临时文件、写了半截
	foreign, err := os.OpenFile(path+".tmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := foreign.WriteString(`{"auth":{"accessToken":"gateway-`); err != nil {
		t.Fatal(err)
	}

	// 面板此刻写入同一账号（正常路径：写临时文件 + rename）
	if _, err := WriteFileAtomic(path, []byte(panelDoc), 0o600); err != nil {
		t.Fatalf("WriteFileAtomic: %v", err)
	}

	// 网关把剩下半截写完（它的句柄仍指向原来那个 inode）
	if _, err := foreign.WriteString(`partial"}}`); err != nil {
		t.Fatal(err)
	}
	if err := foreign.Close(); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != panelDoc {
		t.Errorf("正式凭证被并发写入方污染：\n  got  %q\n  want %q\n"+
			"（面板与网关共用 <path>.tmp；rename 之后网关句柄仍指向正式 inode，其后续字节会写进凭证文件）",
			got, panelDoc)
	}
	if !json.Valid(got) {
		t.Errorf("正式凭证已不是合法 JSON: %q", got)
	}
}

// TestWriteFileAtomicLeavesForeignTempAlone 面板不得吃掉/改写/删除别的进程的临时文件。
// 网关的临时名是固定 <path>.tmp，面板若用它，等于把另一个进程正在写的文件改名带走。
func TestWriteFileAtomicLeavesForeignTempAlone(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "workbuddy-uid-2.json")
	sentinel := "SENTINEL-FOREIGN-TEMP"
	if err := os.WriteFile(path+".tmp", []byte(sentinel), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := WriteFileAtomic(path, []byte(`{"account":{"uid":"uid-2"}}`), 0o600); err != nil {
		t.Fatalf("WriteFileAtomic: %v", err)
	}

	got, err := os.ReadFile(path + ".tmp")
	if err != nil {
		t.Fatalf("并发写入方的临时文件被面板删掉/改名带走了: %v", err)
	}
	if string(got) != sentinel {
		t.Errorf("并发写入方的临时文件被面板改写: %q，期望 %q", got, sentinel)
	}

	// 面板自己的临时文件不能留在目录里
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		base := filepath.Base(path)
		if e.Name() != base && e.Name() != base+".tmp" && strings.Contains(e.Name(), ".tmp") {
			t.Errorf("面板留下了临时文件: %s", e.Name())
		}
	}

	// 两侧的 glob（workbuddy*.json）都不该捞到多余文件
	matches, err := filepath.Glob(filepath.Join(dir, "workbuddy*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 || filepath.Base(matches[0]) != filepath.Base(path) {
		t.Errorf("目录里出现了会被 workbuddy*.json 捞到的多余文件: %v", matches)
	}
}
