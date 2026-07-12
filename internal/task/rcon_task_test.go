package task

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/service"
	"go.etcd.io/bbolt"
)

func openTaskTestDB(t *testing.T) *bbolt.DB {
	t.Helper()
	db, err := bbolt.Open(filepath.Join(t.TempDir(), "test.db"), 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := db.Update(func(tx *bbolt.Tx) error {
		for _, name := range []string{"rcons", "rcon_tasks"} {
			if _, err := tx.CreateBucketIfNotExists([]byte(name)); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestValidateCronExpression(t *testing.T) {
	if err := ValidateCronExpression("*/15 * * * *"); err != nil {
		t.Fatalf("valid cron rejected: %v", err)
	}
	if err := ValidateCronExpression("not a cron"); err == nil {
		t.Fatal("invalid cron was accepted")
	}
}

func TestExecuteRconTaskPersistsSuccess(t *testing.T) {
	db := openTaskTestDB(t)
	if err := service.PutRconCommand(db, "command-1", database.RconCommand{Command: "Broadcast", Remark: "announce"}); err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	rconTask := database.RconTask{
		UUID:       "task-1",
		Name:       "announce",
		RconUUID:   "command-1",
		Content:    "server maintenance",
		Cron:       "0 * * * *",
		Enabled:    true,
		CreatedAt:  now,
		UpdatedAt:  now,
		LastStatus: "never",
	}
	if err := service.AddRconTask(db, rconTask); err != nil {
		t.Fatal(err)
	}

	var executed string
	err := executeRconTask(db, rconTask.UUID, func(command string) (string, error) {
		executed = command
		return "sent", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if executed != "Broadcast server maintenance" {
		t.Fatalf("unexpected command: %q", executed)
	}
	stored, err := service.GetRconTask(db, rconTask.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.LastStatus != "success" || stored.LastResult != "sent" || stored.RunCount != 1 {
		t.Fatalf("unexpected task result: %#v", stored)
	}
}

func TestExecuteRconTaskPersistsFailure(t *testing.T) {
	db := openTaskTestDB(t)
	if err := service.PutRconCommand(db, "command-1", database.RconCommand{Command: "Save"}); err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	rconTask := database.RconTask{
		UUID:       "task-1",
		Name:       "save",
		RconUUID:   "command-1",
		Cron:       "0 * * * *",
		Enabled:    true,
		CreatedAt:  now,
		UpdatedAt:  now,
		LastStatus: "never",
	}
	if err := service.AddRconTask(db, rconTask); err != nil {
		t.Fatal(err)
	}

	expected := errors.New("connection failed")
	err := executeRconTask(db, rconTask.UUID, func(string) (string, error) {
		return "", expected
	})
	if !errors.Is(err, expected) {
		t.Fatalf("expected execution error, got %v", err)
	}
	stored, err := service.GetRconTask(db, rconTask.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.LastStatus != "failed" || stored.LastError != expected.Error() || stored.RunCount != 1 {
		t.Fatalf("unexpected failed result: %#v", stored)
	}
}
