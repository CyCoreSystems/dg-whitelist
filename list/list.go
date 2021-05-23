package list

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

const defaultDBFile = "/var/lib/dg/whitelist.db"

const (
	ListWhite = "white"
	ListGrey = "grey"
	ListBlack = "black"
)

type Item struct {
	// List is the list to which this item belongs: white, grey, or black
	List string `json:"list" form:"list"`
	Address string `json:"address" form:"address"`
	Added time.Time `json:"added" form:"added"`
}

type DB interface {
	Add(list, address string) error
	Remove(list, address string) error
	Get(list string) ([]*Item,error)

	Close() error
}

type boltDB struct {
	db *bbolt.DB
}

func Open(fn string) (DB,error) {
	if fn == "" {
		fn = defaultDBFile
	}

	db, err := bbolt.Open(fn, 0o0600, &bbolt.Options{
		Timeout: time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open db file %q: %w", fn, err)
	}

	// Make sure we have the buckets
	err = db.Update(func(tx *bbolt.Tx) error {
		var err error

		_, err = tx.CreateBucketIfNotExists([]byte(ListBlack))
		if err != nil {
			return fmt.Errorf("failed to create blacklist bucket: %w", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(ListGrey))
		if err != nil {
			return fmt.Errorf("failed to create greylist bucket: %w", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(ListWhite))
		if err != nil {
			return fmt.Errorf("failed to create whitelist bucket: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to ensure buckets exist: %w", err)
	}

	return &boltDB{db}, nil
}

func (db *boltDB) hash(list, addr string) []byte {
	hb := sha256.Sum256([]byte(fmt.Sprintf("%s:%s", list, addr)))

	return hb[:]
}

func (db *boltDB) Close() error {
	return db.db.Close()
}

func (db *boltDB) Add(list, addr string) error {
	return db.db.Update(func(tx *bbolt.Tx) error {
		item := &Item{
			List: list,
			Address: addr,
			Added: time.Now(),
		}

		data, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("failed to marshal item: %w", err)
		}

		b := tx.Bucket([]byte(list))
		if b == nil {
			return fmt.Errorf("bucket %q does not exist", list)
		}

		return tx.Bucket([]byte(list)).Put(db.hash(list, addr), data)
	})
}

func (db *boltDB) Remove(list, addr string) error {
	return db.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(list)).Delete(db.hash(list, addr))
	})
}

func (db *boltDB) Get(list string) ([]*Item,error) {
	var out []*Item

	err := db.db.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket([]byte(list)).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			i := new(Item)

			if err := json.Unmarshal(v, i); err != nil {
				return fmt.Errorf("failed to unmarshal %s: %w", string(k), err)
			}

			out = append(out, i)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}
