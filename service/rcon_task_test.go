package service

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/zaigie/palworld-server-tool/internal/database"
	"go.etcd.io/bbolt"
)

func openRconTaskTestDB(t *testing.T) *bbolt.DB {
	t.Helper()
	db, err := bbolt.Open(filepath.Join(t.TempDir(), "test.db"), 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(rconTasksBucket))
		return err
	}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestRconTaskCRUDAndExecutionResult(t *testing.T) {
	db := openRconTaskTestDB(t)
	createdAt := time.Date(2026, time.July, 12, 10, 0, 0, 0, time.UTC)
	rconTask := database.RconTask{
		UUID:       "task-1",
		Name:       "daily save",
		RconUUID:   "rcon-1",
		Cron:       "0 4 * * *",
		Enabled:    true,
		CreatedAt:  createdAt,
		UpdatedAt:  createdAt,
		LastStatus: "never",
	}

	if err := AddRconTask(db, rconTask); err != nil {
		t.Fatal(err)
	}
	stored, err := GetRconTask(db, rconTask.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Name != rconTask.Name || stored.RconUUID != rconTask.RconUUID {
		t.Fatalf("unexpected stored task: %#v", stored)
	}

	tasks, err := ListRconTasksByCommand(db, rconTask.RconUUID)
	if err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 || tasks[0].UUID != rconTask.UUID {
		t.Fatalf("unexpected command tasks: %#v", tasks)
	}

	ranAt := createdAt.Add(time.Hour)
	if err := UpdateRconTaskExecution(db, rconTask.UUID, "success", "ok", "", ranAt); err != nil {
		t.Fatal(err)
	}
	stored, err = GetRconTask(db, rconTask.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.RunCount != 1 || stored.LastStatus != "success" || stored.LastResult != "ok" {
		t.Fatalf("execution result was not persisted: %#v", stored)
	}
	if stored.LastRunAt == nil || !stored.LastRunAt.Equal(ranAt) {
		t.Fatalf("unexpected last run time: %#v", stored.LastRunAt)
	}

	if err := DeleteRconTask(db, rconTask.UUID); err != nil {
		t.Fatal(err)
	}
	if _, err := GetRconTask(db, rconTask.UUID); err != ErrNoRecord {
		t.Fatalf("expected ErrNoRecord, got %v", err)
	}
}
