package service

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/zaigie/palworld-server-tool/internal/database"
	"go.etcd.io/bbolt"
)

func AddBackup(db *bbolt.DB, backup database.Backup) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("backups"))
		v, err := json.Marshal(backup)
		if err != nil {
			return err
		}
		if err := b.Put([]byte(backup.BackupId), v); err != nil {
			return err
		}
		return nil
	})
}

func GetBackup(db *bbolt.DB, backupId string) (database.Backup, error) {
	var backup database.Backup
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("backups"))
		v := b.Get([]byte(backupId))
		if v == nil {
			return ErrNoRecord
		}
		return json.Unmarshal(v, &backup)
	})
	if err != nil {
		if err == ErrNoRecord {
			return backup, ErrNoRecord
		}
		return backup, err
	}
	return backup, nil
}

func DeleteBackup(db *bbolt.DB, backupId string) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("backups"))
		return b.Delete([]byte(backupId))
	})
}

func ListBackups(db *bbolt.DB, startTime, endTime time.Time) ([]database.Backup, error) {
	backups := make([]database.Backup, 0)
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("backups"))
		return b.ForEach(func(k, v []byte) error {
			var backup database.Backup
			if err := json.Unmarshal(v, &backup); err != nil {
				return err
			}
			// 根据时间筛选
			if (startTime.IsZero() || backup.SaveTime.After(startTime)) &&
				(endTime.IsZero() || backup.SaveTime.Before(endTime)) {
				backups = append(backups, backup)
			}
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].SaveTime.Before(backups[j].SaveTime)
	})
	return backups, nil
}
