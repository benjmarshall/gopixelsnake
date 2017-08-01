package scores

import (
	"strconv"

	"github.com/kniren/gota/dataframe"
)

// Type defines the data structure of the scores object
type Type struct {
	scoresRecord [][]string
}

// NewScores creates a new scores struct
func NewScores() Type {
	t := new(Type)
	t.scoresRecord = [][]string{
		[]string{"Name", "Points"},
	}
	return *t
}

// AddScore pushes a new high score into the scores table
func (t *Type) AddScore(score int, name string) {
	t.scoresRecord = append(t.scoresRecord, []string{name, strconv.Itoa(score)})
	return
}

// GetTopScores returns the top n scores as a string slice, ordered in desecnding order of points
func (t *Type) GetTopScores(n int) [][]string {
	topScoresSlice := [][]string{
		[]string{"Pos.", "Name", "Points"},
	}

	if len(t.scoresRecord) > 1 {
		scoresTable := dataframe.LoadRecords(t.scoresRecord)
		if scoresTable.Err != nil {
			panic(scoresTable.Err)
		}
		orderedScores := scoresTable.Arrange(dataframe.RevSort("Points"))
		if orderedScores.Err != nil {
			panic(orderedScores.Err)
		}
		records := orderedScores.Records()

		for i, record := range records {
			if i == 0 {
				continue
			}
			topScoresSlice = append(topScoresSlice, []string{strconv.Itoa(i), record[0], record[1]})
		}
	}

	return topScoresSlice
}
