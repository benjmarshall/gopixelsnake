package scores

import (
	"encoding/csv"
	"strconv"

	"github.com/kniren/gota/dataframe"
	"github.com/shibukawa/configdir"
)

// Type defines the data structure of the scores object
type Type struct {
	scoresRecord [][]string
	scoresFile   string
	configDirs   configdir.ConfigDir
	numScores    int
}

// NewScores creates a new scores struct
func NewScores(filename string, numScores int) Type {
	t := new(Type)
	t.scoresFile = filename
	t.scoresRecord = [][]string{
		[]string{"Name", "Points"},
	}
	t.configDirs = configdir.New("benjmarshall", "gopixelsnake")
	t.numScores = numScores
	t.LoadScores()
	return *t
}

// AddScore pushes a new high score into the scores table
func (t *Type) AddScore(score int, name string) {
	record := t.getOrderedScores()
	if len(record) > t.numScores {
		t.scoresRecord = append(record[0:len(record)-1], []string{name, strconv.Itoa(score)})
	} else {
		t.scoresRecord = append(t.scoresRecord, []string{name, strconv.Itoa(score)})
	}
	t.SaveScores()
	return
}

// GetTopScores returns the top n scores as a string slice, ordered in desecnding order of points
func (t *Type) GetTopScores() [][]string {
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

func (t *Type) getOrderedScores() [][]string {

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
		return records
	}

	return t.scoresRecord
}

// GetBottomScore returns the points value of the nth high score in descending order
func (t *Type) GetBottomScore() int {
	scores := t.GetTopScores()
	if len(scores) < t.numScores {
		return 0
	}
	score, err := strconv.Atoi(scores[t.numScores][2])
	if err != nil {
		return 0
	}
	return score
}

// SaveScores saves the scores to a csv file
func (t *Type) SaveScores() {
	// If we haven't got any scores stop now
	if len(t.scoresRecord) < 2 {
		return
	}

	folders := t.configDirs.QueryFolders(configdir.Global)

	f, err := folders[0].Create(t.scoresFile)
	if err != nil {
		return
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	scores := t.scoresRecord[1 : len(t.scoresRecord)+1]

	for _, value := range scores {
		err := w.Write(value)
		if err != nil {
			return
		}
	}
}

// LoadScores loads saved scores from a csv file
func (t *Type) LoadScores() {
	folder := t.configDirs.QueryFolderContainsFile(t.scoresFile)
	if folder != nil {
		f, err := folder.Open(t.scoresFile)
		if err != nil {
			return
		}
		defer f.Close()

		r := csv.NewReader(f)
		records, err := r.ReadAll()
		if err != nil {
			return
		}

		for _, record := range records {
			t.scoresRecord = append(t.scoresRecord, record)
		}
	}

	return
}
