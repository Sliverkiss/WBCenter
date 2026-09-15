package control

import "testing"

func TestLoadConfigReadsOperationalEnvironment(t *testing.T) {
	t.Setenv("WBCC_PASSWORD", "test-password")
	t.Setenv("WBCC_AUTH_DIR", t.TempDir())
	t.Setenv("WBCC_DATA_DIR", t.TempDir())
	t.Setenv("WBCC_READ_ONLY", "true")
	t.Setenv("WBCC_TIMEZONE", "UTC")
	t.Setenv("WBCC_TIMEOUT_SECONDS", "45")
	cfg, err := LoadConfig(t.TempDir() + "/missing.json")
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.ReadOnly || cfg.Timezone != "UTC" || cfg.TimeoutSeconds != 45 {
		t.Fatalf("environment overrides were not applied: %#v", cfg)
	}
}

func TestLoadConfigRejectsInvalidOperationalEnvironment(t *testing.T) {
	t.Setenv("WBCC_PASSWORD", "test-password")
	t.Setenv("WBCC_AUTH_DIR", t.TempDir())
	t.Setenv("WBCC_READ_ONLY", "sometimes")
	if _, err := LoadConfig(t.TempDir() + "/missing.json"); err == nil {
		t.Fatal("invalid WBCC_READ_ONLY was accepted")
	}
}
