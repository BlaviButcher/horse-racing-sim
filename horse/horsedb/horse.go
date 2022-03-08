package horsedb

import (
	"github.com/rs/xid"
	"github.com/vmihailenco/msgpack/v5"
)

type Horse struct {
	ID   xid.ID
	name string
	mmr  int
}

func NewHorse(name string) *Horse {
	return &Horse{
		ID:   xid.New(),
		name: name,
		mmr:  1500,
	}
}

// Marshal converts Horse struct to a byte array and returns it.
func (h *Horse) Marshal() ([]byte, error) {
	return msgpack.Marshal(h)
}

// UnMarshal converts a byte array and returns an error.
func (h *Horse) Unmarshal(data []byte) error {
	return msgpack.Unmarshal(data, h)
}

// Updates ID based on data given (marshalled struct)
func (h *Horse) updateWithKey(data []byte) {
	k := HorseKeyFromBytes(data)
	h.ID = k.HorseID
}

func (h *Horse) Key() Key           { return h.HorseKey() }
func (h *Horse) HorseKey() HorseKey { return NewHorseKey(h.ID) }
