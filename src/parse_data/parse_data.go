package parse_data

import (
	"encoding/csv"
	"github.com/fxtlabs/date"
	"os"
	"strings"
)

const DateLayout = "2006-01-02"

type Matchup struct {
	Country1 string
	Country2 string
}

type Result struct {
	Matchup
	Date          date.Date
	Country1Score string
	Country2Score string
	Tournament    string
	Neutral       string
}

type Ranking struct {
	Country        string
	Rank           string
	TotalPoints    string
	PreviousPoints string
	RankChange     string
	Confederation  string
	Date           date.Date
}

type HistoricalRatios struct {
	Matchup
	Games        string
	Country1Win  string
	Country1Loss string
	Country1Draw string
}

type FullMatchData map[string]string

type knockoutStageData struct {
	FullMatchData
	Shootout
}

type Shootout struct {
	Date date.Date
	Matchup
	Winner string
}

func getHistoricalData() (r []Result) {
	csvData, err := ReadCsVFile("data/historical-results.csv")
	if err != nil {
		panic(err)
	}
	var results []Result

	for i, line := range csvData {
		if i != 0 {
			matchDate, _ := date.Parse(DateLayout, line[0])
			if err != nil {
				panic(err)
			}
			data := Result{
				Matchup: Matchup{
					Country1: line[1],
					Country2: line[2],
				},
				Date:          matchDate,
				Country1Score: line[3],
				Country2Score: line[4],
				Tournament:    line[5],
				Neutral:       line[8],
			}

			results = append(results, data)
		}

	}
	return results
}

func getRankings() (r []Ranking) {
	csvData, err := ReadCsVFile("data/ranking.csv")
	if err != nil {
		panic(err)
	}
	var rankings []Ranking

	for i, line := range csvData {
		if i != 0 {
			rankDate, _ := date.Parse(DateLayout, line[7])
			if err != nil {
				panic(err)
			}
			data := Ranking{
				Country:        line[1],
				Rank:           line[0],
				TotalPoints:    line[3],
				PreviousPoints: line[4],
				RankChange:     line[5],
				Confederation:  line[6],
				Date:           rankDate,
			}

			rankings = append(rankings, data)
		}

	}
	return rankings
}

func getRatios() (r []HistoricalRatios) {
	csvData, err := ReadCsVFile("data/historical_win-loose-draw_ratios.csv")
	if err != nil {
		panic(err)
	}
	var ratios []HistoricalRatios

	for i, line := range csvData {
		if i != 0 {
			data := HistoricalRatios{
				Matchup: Matchup{
					Country1: line[0],
					Country2: line[1],
				},
				Games:        line[2],
				Country1Win:  line[3],
				Country1Loss: line[4],
				Country1Draw: line[5],
			}

			ratios = append(ratios, data)
		}

	}
	return ratios
}

func getShootouts() (r []Shootout) {
	csvData, err := ReadCsVFile("data/shootouts.csv")
	if err != nil {
		panic(err)
	}
	var shootout []Shootout

	for i, line := range csvData {
		if i != 0 {
			shootoutDate, _ := date.Parse(DateLayout, line[0])
			if err != nil {
				panic(err)
			}
			data := Shootout{
				Date: shootoutDate,
				Matchup: Matchup{
					Country1: line[1],
					Country2: line[2],
				},
				Winner: line[3],
			}

			shootout = append(shootout, data)
		}

	}
	return shootout
}

func GetScheduledMatches() (m []Result) {
	csvData, err := ReadCsVFile("data/matches-schedule.csv")
	if err != nil {
		panic(err)
	}
	var match []Result

	for i, line := range csvData {
		if i != 0 {
			dateSlice := strings.Split(line[1], "/")
			matchDate, _ := date.Parse(DateLayout, dateSlice[2]+"-"+dateSlice[1]+"-"+dateSlice[0])
			if err != nil {
				panic(err)
			}
			data := Result{
				Matchup: Matchup{
					Country1: line[2],
					Country2: line[3],
				},
				Country1Score: "",
				Country2Score: "",
				Tournament:    "FIFA World Cup",
				Date:          matchDate,
				Neutral:       "TRUE",
			}
			match = append(match, data)
		}
	}
	return match
}

func ReadCsVFile(filename string) ([][]string, error) {
	// Open CSV file
	fileContent, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}

	defer func(fileContent *os.File) {
		err := fileContent.Close()
		if err != nil {

		}
	}(fileContent)

	// Read File into a Variable
	lines, err := csv.NewReader(fileContent).ReadAll()

	return lines, err
}

func GetAllData() (hr []Result, ranks []Ranking, ratio []HistoricalRatios, s []Shootout, m []Result) {
	historicalResults := getHistoricalData()
	rankings := getRankings()
	ratios := getRatios()
	shootouts := getShootouts()
	matches := GetScheduledMatches()

	return historicalResults, rankings, ratios, shootouts, matches
}
