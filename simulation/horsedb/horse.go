package horsedb

import (
	"math/rand"

	"github.com/rs/xid"
	"github.com/vmihailenco/msgpack/v5"
)

const HistoryMMRKept = 50
const MinMMR = 1000
const MaxMMR = 3000
const Variance = 500
const Movement = 3

type Horse struct {
	ID             xid.ID
	Name           string
	MMR            []int
	AvgMMR         int
	RawMMR         int // Moves around based on previous race mmr + new change in mmr
	MMRChange      int
	WinProbability float64
	GeneticMMR     int
	MMRVariance    int
}

func NewHorse(name string) *Horse {

	startingMMR := rand.Intn(MaxMMR-MinMMR) + MinMMR

	MMR := make([]int, HistoryMMRKept)
	for i := 0; i < HistoryMMRKept; i++ {
		MMR[i] = startingMMR
	}

	return &Horse{
		ID:             xid.New(),
		Name:           name,
		MMR:            MMR,
		AvgMMR:         startingMMR,
		RawMMR:         startingMMR,
		WinProbability: 0,
		MMRVariance:    Variance,
		GeneticMMR:     startingMMR,
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
