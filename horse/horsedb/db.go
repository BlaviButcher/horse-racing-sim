package horsedb

import (
	"fmt"
	"sync"

	"github.com/dgraph-io/badger/v3"
)

type DB struct {
	DB *badger.DB

	writeMu *WriteMu

	// only close the database once
	closeOnce sync.Once
}

// Set element in the database
// Not efficient for lots of entries need a defined function for that
func (db DB) set(items ...Value) error {

	return db.DB.Update(func(txn *badger.Txn) error {
		for _, i := range items {
			if err := setValue(txn, i); err != nil {
				return fmt.Errorf("writing: %w", err)
			}
		}
		return nil
	})
}

// Sets a new entry in the database
func setValue(txn *badger.Txn, vs ...Value) error {

	for _, v := range vs {
		data, err := v.Marshal()
		if err != nil {
			return fmt.Errorf("marshalling: %w", err)
		}
		e := badger.NewEntry(makeKey(v.Key()), data)

		err = txn.SetEntry(e)
		if err != nil {
			return fmt.Errorf("setting to txn: %w", err)
		}
	}
	return nil
}

func makeKey(k Key) []byte { return k.Key(make([]byte, 0, k.KeySize())) }

func fromItem(item *badger.Item, v Value) error {
	err := item.Value(func(val []byte) error {
		return v.Unmarshal(val)
	})
	if err != nil {
		return fmt.Errorf("unmarshalling: %w", err)
	}
	v.updateWithKey(item.Key())

	return nil
}

// OpenDatabase creates a new database object
// Returns database object and an error
func OpenDatabase(path string) (*DB, error) {
	opts := badger.DefaultOptions(path)
	opts.Logger = nil

	bdb, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("opening badger database: %w", err)
	}

	db := &DB{
		DB:      bdb,
		writeMu: new(WriteMu),
	}

	return db, nil
}

// Close the database properly
// Should always be called, at least in a simple defer
// but ideally in a context where the error can be checked
func (db *DB) Close() error {
	var err error
	db.closeOnce.Do(func() {
		db.writeMu.StartStopWorldWrites()
		defer db.writeMu.DoneStopWorldWrites()

		if err = db.DB.Close(); err != nil {
			err = fmt.Errorf("closing database: %w", err)
			return
		}
	})
	return err
}
