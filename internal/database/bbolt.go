package database

import (
	"sync"
	"time"

	"github.com/zaigie/palworld-server-tool/internal/logger"
	"go.etcd.io/bbolt"
	"github.com/zaigie/palworld-server-tool/internal/config"
)

var db *bbolt.DB
var once sync.Once

func InitDB() *bbolt.DB {
	db_, err := bbolt.Open(config.Config.Db.path, 0600, &bbolt.Options{Timeout: 1 * time.Minute})
	if err != nil {
		logger.Panic(err)
	}
	// players
	err = db_.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("players"))
		return err
	})
	if err != nil {
		logger.Panic(err)
	}
	// guilds
	err = db_.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("guilds"))
		return err
	})
	if err != nil {
		logger.Panic(err)
	}
	// rcons
	err = db_.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("rcons"))
		return err
	})
	if err != nil {
		logger.Panic(err)
	}
	// backups
	err = db_.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("backups"))
		return err
	})
	if err != nil {
		logger.Panic(err)
	}
	return db_
}

func GetDB() *bbolt.DB {
	once.Do(func() {
		db = InitDB()
	})
	return db
}
