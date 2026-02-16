package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.IdeasDirectory == "" {
		t.Error("IdeasDirectory should not be empty")
	}
	if cfg.Editor != "vim" {
		t.Errorf("Editor: got %q, want %q", cfg.Editor, "vim")
	}
}

func TestLoad_FromFile(t *testing.T) {
	dir := t.TempDir()
	ideasDir := filepath.Join(dir, "ideas")
	os.Mkdir(ideasDir, 0755)

	configContent := `ideas_directory = "` + ideasDir + `"
editor = "nano"
`
	configPath := filepath.Join(dir, "config.toml")
	os.WriteFile(configPath, []byte(configContent), 0644)

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.IdeasDirectory != ideasDir {
		t.Errorf("IdeasDirectory: got %q, want %q", cfg.IdeasDirectory, ideasDir)
	}
	if cfg.Editor != "nano" {
		t.Errorf("Editor: got %q, want %q", cfg.Editor, "nano")
	}
}

func TestLoad_NonexistentFile(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.toml")
	if err != nil {
		t.Fatalf("Load should return defaults for nonexistent file: %v", err)
	}

	if cfg.IdeasDirectory == "" {
		t.Error("should return default config")
	}
}

func TestValidate_MissingDirectory(t *testing.T) {
	cfg := &Config{
		IdeasDirectory: "/nonexistent/path/ideas",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for nonexistent directory")
	}
}

func TestValidate_EmptyDirectory(t *testing.T) {
	cfg := &Config{
		IdeasDirectory: "",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for empty directory")
	}
}

func TestValidate_ValidDirectory(t *testing.T) {
	dir := t.TempDir()
	cfg := &Config{
		IdeasDirectory: dir,
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for valid directory: %v", err)
	}
}

func TestValidate_NotADirectory(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "notadir")
	os.WriteFile(filePath, []byte("test"), 0644)

	cfg := &Config{
		IdeasDirectory: filePath,
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("expected error when ideas_directory is a file")
	}
}

func TestExpandHome(t *testing.T) {
	homeDir, _ := os.UserHomeDir()

	tests := []struct {
		input string
		want  string
	}{
		{"~/ideas", filepath.Join(homeDir, "ideas")},
		{"~/", homeDir},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
		{"", ""},
	}

	for _, tt := range tests {
		got := expandHome(tt.input)
		if got != tt.want {
			t.Errorf("expandHome(%q): got %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestConfigPath(t *testing.T) {
	path := ConfigPath()
	if path == "" {
		t.Error("ConfigPath should return a non-empty path")
	}
}
