package field

import (
	"fmt"

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
	f.postRaceMMRAdjustment()
}

func (f *Field) caclulateHorseWinProbability() {
	for _, h := range f.horses {

		h.WinProbability = float64(h.AvgMMR) / float64(f.totalMMR)
	}
}

func (f *Field) setTotalMMR() {
	for _, h := range f.horses {
		f.totalMMR += h.AvgMMR
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
	f.totalMMR -= winner.AvgMMR
	// remove horse from the line up so we can find the next place
	delete(f.horses, winner.ID)

}

func (f *Field) postRaceMMRAdjustment() {
	for i, placedHorse := range f.RaceResults {
		// total MMR of everyone who lost to horse but had better avgMMR
		totalSuperiorToMMR := 0
		// total MMR of everyone who beat horse and had better avgMMR
		totalInferiorToMMR := 0

		// TODO: The bug is here due to the fact that we are making changes to the horses pointers before we are done with the for loop
		// maybe create a field to hold the change until we are done here
		for j, opponentHorse := range f.RaceResults {
			if j < i && placedHorse.AvgMMR < opponentHorse.AvgMMR {
				fmt.Println(placedHorse.AvgMMR)
				fmt.Println(opponentHorse.AvgMMR)
				// should be adding the diff of placedHorse avg and oppenentHorse avg
				totalSuperiorToMMR += opponentHorse.AvgMMR
			} else if j > i && placedHorse.AvgMMR > opponentHorse.AvgMMR {
				totalInferiorToMMR += opponentHorse.AvgMMR
			}
		}
		fmt.Println(totalSuperiorToMMR)
		// fmt.Println(totalInferiorToMMR)

		// change in MMR plus total horse beaten
		mmrChange := (totalSuperiorToMMR - totalInferiorToMMR) + i

		// slide slice up 1
		newMMR := placedHorse.MMR[1:]
		newMMR = append(newMMR, placedHorse.RawMMR+mmrChange)

		// get new average MMR
		totalMMR := 0
		for _, mmr := range newMMR {
			totalMMR += mmr
		}
		placedHorse.AvgMMR = totalMMR / len(newMMR)

		placedHorse.MMR = newMMR

	}
}
