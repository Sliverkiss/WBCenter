package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

// TestOverwriteInPlaceTightensExistingPerm 回退路径（Docker 单文件挂载，rename 不可用）
// 必须把结果权限收敛到调用方要求的值。
//
// 为什么：os.OpenFile/os.WriteFile 的 perm 只对「新建」生效，已存在的文件会被静默忽略。
// 而回退路径的目标文件通常是宿主机上人工创建或按 umask 建出来的（config.json / 凭证文件），
// mode 往往比 0600 宽——于是面板声称「原子写 0600」却把一个组可读/全局可读的凭证留在盘上。
func TestOverwriteInPlaceTightensExistingPerm(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "mounted.json")
	if err := os.WriteFile(path, []byte("original-longer-content"), 0o600); err != nil {
		t.Fatal(err)
	}
	// 宿主机侧就是宽权限（人工创建 / umask 决定）
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := overwriteInPlace(path, []byte("short"), 0o600); err != nil {
		t.Fatalf("overwriteInPlace: %v", err)
	}

	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := fi.Mode().Perm(); got != 0o600 {
		t.Errorf("回退路径没有收敛权限：mode = %#o，期望 0600"+
			"（凭证/配置含敏感信息，同宿主或同挂载的其它用户不应可读）", got)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "short" {
		t.Errorf("内容 = %q，期望 %q（截断仍须彻底）", got, "short")
	}
}

// TestOverwriteInPlaceCreatesWithPerm 新建分支的结果权限也必须等于 perm，
// 不受 umask 影响；行为不得回退。
func TestOverwriteInPlaceCreatesWithPerm(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fresh.json")
	if err := overwriteInPlace(path, []byte("x"), 0o600); err != nil {
		t.Fatalf("overwriteInPlace: %v", err)
	}
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := fi.Mode().Perm(); got != 0o600 {
		t.Errorf("新建文件 mode = %#o，期望 0600", got)
	}
}
