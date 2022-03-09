package field

import (
	"github.com/blavi/horse/simulation/horsedb"
	"github.com/jmcvetta/randutil"
	"github.com/rs/xid"
)

type Field struct {
	horses      map[xid.ID]*horsedb.Horse
	totalMMR    int
	RaceResults []*horsedb.Horse
}

func NewField(horses []*horsedb.Horse) *Field {
	horseMap := make(map[xid.ID]*horsedb.Horse)

	for _, h := range horses {
		horseMap[h.ID] = h
	}

	f := &Field{
		horses:      horseMap,
		totalMMR:    0,
		RaceResults: []*horsedb.Horse{},
	}
	f.setTotalMMR()
	return f
}

func (f *Field) Race() {

	// length of horses will change throughout loop, so need static
	fieldSize := len(f.horses)
	for i := 0; i < fieldSize; i++ {
		f.caclulateHorseWinProbability()
		f.getNextPlace()

		// fmt.Println(f.RaceResults[len(f.RaceResults)-1].Name)
	}
}

func (f *Field) caclulateHorseWinProbability() {
	for _, h := range f.horses {
		h.WinProbability = float64(h.MMR) / float64(f.totalMMR)
	}
}

func (f *Field) setTotalMMR() {
	for _, h := range f.horses {
		f.totalMMR += h.MMR
	}
}

// gets the next winner of the remaining horses in field
func (f *Field) getNextPlace() {
	choices := []randutil.Choice{}
	for _, h := range f.horses {
		// converting to "4 siginificant figures" when truncated as int
		choices = append(choices, randutil.Choice{int(h.WinProbability * 10000), h.ID})
	}

	choice, _ := randutil.WeightedChoice(choices)
	id := choice.Item.(xid.ID)

	winner := f.horses[id]
	f.RaceResults = append(f.RaceResults, winner)

	// Remove horse and update MMR
	f.totalMMR -= winner.MMR
	// remove horse from the line up so we can find the next place
	delete(f.horses, winner.ID)

}
