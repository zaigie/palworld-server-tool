package service

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/zaigie/palworld-server-tool/internal/database"
	"go.etcd.io/bbolt"
)

const rconTasksBucket = "rcon_tasks"

func AddRconTask(db *bbolt.DB, task database.RconTask) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(rconTasksBucket))
		if err != nil {
			return err
		}
		value, err := json.Marshal(task)
		if err != nil {
			return err
		}
		return b.Put([]byte(task.UUID), value)
	})
}

func PutRconTask(db *bbolt.DB, task database.RconTask) error {
	return AddRconTask(db, task)
}

func GetRconTask(db *bbolt.DB, taskUUID string) (database.RconTask, error) {
	var task database.RconTask
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(rconTasksBucket))
		if b == nil {
			return ErrNoRecord
		}
		value := b.Get([]byte(taskUUID))
		if value == nil {
			return ErrNoRecord
		}
		return json.Unmarshal(value, &task)
	})
	return task, err
}

func ListRconTasks(db *bbolt.DB) ([]database.RconTask, error) {
	tasks := make([]database.RconTask, 0)
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(rconTasksBucket))
		if b == nil {
			return nil
		}
		return b.ForEach(func(_, value []byte) error {
			var task database.RconTask
			if err := json.Unmarshal(value, &task); err != nil {
				return err
			}
			tasks = append(tasks, task)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})
	return tasks, nil
}

func ListRconTasksByCommand(db *bbolt.DB, rconUUID string) ([]database.RconTask, error) {
	tasks, err := ListRconTasks(db)
	if err != nil {
		return nil, err
	}
	matched := make([]database.RconTask, 0)
	for _, task := range tasks {
		if task.RconUUID == rconUUID {
			matched = append(matched, task)
		}
	}
	return matched, nil
}

func DeleteRconTask(db *bbolt.DB, taskUUID string) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(rconTasksBucket))
		if b == nil {
			return nil
		}
		return b.Delete([]byte(taskUUID))
	})
}

func UpdateRconTaskExecution(db *bbolt.DB, taskUUID, status, result, runError string, ranAt time.Time) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(rconTasksBucket))
		if b == nil {
			return ErrNoRecord
		}
		value := b.Get([]byte(taskUUID))
		if value == nil {
			return ErrNoRecord
		}
		var task database.RconTask
		if err := json.Unmarshal(value, &task); err != nil {
			return err
		}
		task.LastRunAt = &ranAt
		task.LastStatus = status
		task.LastResult = result
		task.LastError = runError
		task.RunCount++
		task.UpdatedAt = ranAt
		updated, err := json.Marshal(task)
		if err != nil {
			return err
		}
		return b.Put([]byte(taskUUID), updated)
	})
}
