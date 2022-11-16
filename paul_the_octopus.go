package main

import (
	"net/http"
	"os"
	"preprocessing"
)

func main() {
	//file, err := parse_data.ReadCsVFile("data/historical_data_cleaned.csv")
	//if err != nil {
	//	return
	//}
	//
	//for _, i := range file {
	//	fmt.Println(i)
	//}
	preprocessing.PreProcessData(os.Stdout)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//if r.URL.Path != "/" {
	//	http.NotFound(w, r)
	//	return
	//}

}
