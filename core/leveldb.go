package core

import (
	"github.com/syndtr/goleveldb/leveldb"
)

var dbPath = "/var/lib/trojan-manager"

// GetValue Get Leveldb value
func GetValue(key string) (string, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return "", err
	}
	defer db.Close()
	result, err := db.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// SetValue Set the leveldb value
func SetValue(key string, value string) error {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Put([]byte(key), []byte(value), nil)
}

// DelValue delete
func DelValue(key string) error {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Delete([]byte(key), nil)
}
