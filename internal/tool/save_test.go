package tool

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/zaigie/palworld-server-tool/internal/config"
)

func TestBackupHonorsConfiguredSaveSourceMode(t *testing.T) {
	store, err := config.Open(filepath.Join(t.TempDir(), "config.db"))
	if err != nil {
		t.Fatalf("open config store: %v", err)
	}
	defer store.Close()
	config.SetCurrent(store)

	settings := store.Config()
	settings.Save.SourceMode = "directory"
	settings.Save.Path = "http://game-server:8081/sync"
	if err := store.Update(settings, ""); err != nil {
		t.Fatalf("save directory settings: %v", err)
	}
	if _, err := Backup(); err == nil || !strings.Contains(err.Error(), "directory source") {
		t.Fatalf("directory mode with URL error = %v, want directory source validation", err)
	}

	settings.Save.SourceMode = "agent"
	settings.Save.Path = t.TempDir()
	if err := store.Update(settings, ""); err != nil {
		t.Fatalf("save agent settings: %v", err)
	}
	if _, err := Backup(); err == nil || !strings.Contains(err.Error(), "agent source") {
		t.Fatalf("agent mode with local path error = %v, want agent source validation", err)
	}
}
