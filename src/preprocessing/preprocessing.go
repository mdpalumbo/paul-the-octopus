package preprocessing

import (
	"encoding/csv"
	"github.com/fxtlabs/date"
	"io"
	"log"
	"math"
	"os"
	"parse_data"
	"sort"
)

type result = parse_data.Result
type ranking = parse_data.Ranking
type fullData = parse_data.FullMatchData
type ratio = parse_data.HistoricalRatios
type rankingsByTeam map[string][]ranking

type finalDataSet = []fullData

func PreProcessData(w io.Writer) {
	results, ranks, ratios, _, _ := parse_data.GetAllData()
	//results = append(results, scheduledMatches...)
	finalDataMap := attachDataToResults(w, results, ranks, ratios)
	writeMapToCSV(finalDataMap)
}

func writeMapToCSV(dataSet finalDataSet) {
	var sliceForCSV [][]string
	var keys []string
	// create array of strings in array  for csv writer
	for i, data := range dataSet {

		var valueArray []string
		for k := range data {
			if i == 0 {
				keys = append(keys, k)
			}
		}
		sort.Sort(sort.StringSlice(keys))
		if i == 0 {
			sliceForCSV = append(sliceForCSV, keys)
		}
		for _, k := range keys {
			val := data[k]
			valueArray = append(valueArray, val)
		}
		sliceForCSV = append(sliceForCSV, valueArray)
	}

	f, err := os.Create("data/historical_data_cleaned.csv")
	defer f.Close()

	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range sliceForCSV {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
}

// fd fullData
func attachDataToResults(w io.Writer, results []result, rankings []ranking, ratios []ratio) (f finalDataSet) {
	ranksByTeam := groupRankingsByTeam(w, rankings)
	var fullTrainingData finalDataSet
	for _, result := range results {
		var resultData = make(fullData)
		r1, r2 := getLatestRankingsForResult(w, result, ranksByTeam)
		ratio := getRatiosForResult(w, result, ratios)
		resultData["Country1"] = result.Country1
		resultData["Country2"] = result.Country2
		resultData["MatchDate"] = result.Date.Format(parse_data.DateLayout)
		resultData["Country1Score"] = result.Country1Score
		resultData["Country2Score"] = result.Country2Score
		resultData["Tournament"] = result.Tournament
		resultData["Neutral"] = result.Neutral
		resultData["HistoricalGameCount"] = ratio.Games
		resultData["Country1WinRatio"] = ratio.Country1Win
		resultData["Country1LossRatio"] = ratio.Country1Loss
		resultData["Country1DrawRatio"] = ratio.Country1Draw
		resultData["Country1Rank"] = r1.Rank
		resultData["Country1RankTotalPoints"] = r1.TotalPoints
		resultData["Country1RankPreviousPoints"] = r1.PreviousPoints
		resultData["Country1RankChange"] = r1.RankChange
		resultData["Country1RankConfederation"] = r1.Confederation
		resultData["Country2Rank"] = r2.Rank
		resultData["Country2RankTotalPoints"] = r2.TotalPoints
		resultData["Country2RankPreviousPoints"] = r2.PreviousPoints
		resultData["Country2RankChange"] = r2.RankChange
		resultData["Country2RankConfederation"] = r2.Confederation
		fullTrainingData = append(fullTrainingData, resultData)
	}
	return fullTrainingData
}

func getRatiosForResult(w io.Writer, res result, ratios []ratio) (r ratio) {
	var wantedRatio = ratio{
		Matchup: parse_data.Matchup{
			Country1: res.Country1,
			Country2: res.Country2,
		},
	}
	for _, _ratio := range ratios {
		if _ratio.Country1 == res.Country1 && _ratio.Country2 == res.Country2 {
			wantedRatio = _ratio
		}
	}
	//_, _ = fmt.Fprintf(w, "Ratio: \t\t%v\n\n", wantedRatio)
	return wantedRatio
}

func getLatestRankingsForResult(w io.Writer, result result, ranksByTeam rankingsByTeam) (r1 ranking, r2 ranking) {
	var rankingsCountry1 = ranksByTeam[result.Country1]
	var rankingsCountry2 = ranksByTeam[result.Country2]

	country1ClosestRank := getClosestRankDate(rankingsCountry1, result)
	country2ClosestRank := getClosestRankDate(rankingsCountry2, result)

	//_, err := fmt.Fprintf(w, "%v\nC1Rank:\t\t%v\nC2Rank:\t\t%v\n", result, country1ClosestRank, country2ClosestRank)
	//if err != nil {
	//	return
	//}
	return country1ClosestRank, country2ClosestRank
}

func getClosestRankDate(rankings []ranking, result result) ranking {
	var closestRankingDate = getClosestDate(rankings, result.Date)
	var countryRank ranking
	for _, rank := range rankings {
		if rank.Date == closestRankingDate {
			countryRank = rank
		}
	}
	return countryRank
}

func getClosestDate(countryRankings []ranking, resultDate date.Date) date.Date {
	minDiff := -1.0
	var minDate date.Date
	for _, rank := range countryRankings {
		diff := math.Abs(float64(resultDate.Sub(rank.Date)))
		if minDiff == -1 || diff < minDiff {
			minDiff = diff
			minDate = rank.Date
		}
	}
	return minDate
}

func groupRankingsByTeam(w io.Writer, rankings []ranking) (rbt rankingsByTeam) {
	ranksByTeam := make(rankingsByTeam)
	for i, rank := range rankings {
		if i != 0 {
			ranksByTeam[rank.Country] = append(ranksByTeam[rank.Country], rank)
		}
	}
	// sort ranks by date
	for _, ranks := range ranksByTeam {
		sort.Slice(ranks, func(i, j int) bool {
			return ranks[i].Date.Before(ranks[j].Date)
		})
	}
	return ranksByTeam
}
