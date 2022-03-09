package simulation

import (
	"fmt"
	"math/rand"

	"github.com/blavi/horse/simulation/field"
	"github.com/blavi/horse/simulation/horsedb"
	"github.com/rs/xid"
)

// Algorithm to pick 14 horses for race
// Takes in horses and returns horses for race
// TODO: implement way in which horses with less races have more chance to race
func getNextRaceGroup(horses []*horsedb.Horse) []*horsedb.Horse {
	if len(horses) < 15 {
		return horses
	}

	horseMap := make(map[xid.ID]*horsedb.Horse)
	horseCount := 0
	for horseCount < 14 {
		h := horses[rand.Intn(len(horses))]
		if horseMap[h.ID] != nil {
			continue
		}

		horseMap[h.ID] = h

	}

	out := make([]*horsedb.Horse, 0, len(horseMap))
	for _, v := range horseMap {
		out = append(out, v)
	}

	return out
}

//TODO: probably need an object of sorts that will store information and hold the current pool of horses also, that we can call from. Rather
// than passing horse addresses around
func SimulateFreshHorses(horses []*horsedb.Horse, races int) ([]*horsedb.Horse, error) {

	// check all horses are new
	for _, h := range horses {
		if h.AvgMMR != 1500 {
			return nil, fmt.Errorf("horse of id %s not ")
		}
	}

	for i := 0; i < races; i++ {
		raceField := field.NewField(getNextRaceGroup(horses))
		raceField.Race()
	}

	return horses, nil

}
