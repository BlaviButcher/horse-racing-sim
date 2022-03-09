package horsedb

import (
	"github.com/dgraph-io/badger/v3"
)

func (db *DB) SetHorses(horses ...*Horse) error {

	if len(horses) == 0 {
		return nil
	} else if len(horses) == 1 {
		return db.set(horses[0])
	}

	txn := db.DB.NewTransaction(true)
	defer txn.Discard()
	for _, h := range horses {
		if err := setValue(txn, h); err != nil {
			return err
		}
	}

	return txn.Commit()

}

func (db *DB) ListHorses() ([]*Horse, error) {

	out := make([]*Horse, 0)

	db.DB.View(func(txn *badger.Txn) error {

		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte{HorseKeyPrefix}
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			m := new(Horse)

			if err := item.Value(func(val []byte) error {
				return m.Unmarshal(val)
			}); err != nil {
				continue
			}
			out = append(out, m)
		}
		return nil
	})

	return out, nil
}
