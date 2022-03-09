package horsedb

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/rs/xid"
)

const (
	HorseKeyPrefix byte = iota
)

type Key interface {
	// Serialises the key
	Key([]byte) []byte

	// Returns size of key when serialised
	KeySize() int
}

// Methods for HorseKey
type HorseKey struct {
	HorseID xid.ID
}

// returns key as HorseKey type
func NewHorseKey(k xid.ID) HorseKey {
	return HorseKey{
		HorseID: k,
	}
}

// Makes HorseKey into horsekey with correct prefix
func (k HorseKey) Key(to []byte) []byte {
	to = append(to, HorseKeyPrefix)
	to = append(to, k.HorseID[:]...)
	return to
}

func (k HorseKey) KeySize() int {
	return 1 + len(k.HorseID)
}

func HorseKeyFromBytes(data []byte) HorseKey {
	id, _ := xid.FromBytes(data[1:])
	return HorseKey{
		HorseID: id,
	}
}

func (k HorseKey) ValueFrom(txn *badger.Txn) (*Horse, error) {
	item, err := txn.Get(k.Key(make([]byte, 0, k.KeySize())))
	if err != nil {
		return nil, err
	}
	v := new(Horse)
	err = item.Value(func(data []byte) error { return v.Unmarshal(data) })
	if err != nil {
		return nil, err
	}
	v.updateWithKey(item.Key())
	return v, nil
}

func (k HorseKey) PopulateValue(txn *badger.Txn, v *Horse) error {
	item, err := txn.Get(k.Key(make([]byte, 0, k.KeySize())))
	if err != nil {
		return err
	}
	err = item.Value(func(data []byte) error { return v.Unmarshal(data) })
	if err != nil {
		return err
	}
	v.updateWithKey(item.Key())
	return nil
}
