package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	cf "polygon/config"
	models "polygon/models"
)

func parseAggregateRequest(c cf.Secrets, ar models.AggregateRequest) string {
	return fmt.Sprintf("https://api.polygon.io/v2/aggs/ticker/%s/range/%d/%s/%s/%s?adjusted=%t&sort=asc&limit=%d&apiKey=%s",
		ar.Ticker, ar.RangeLength, ar.RangeType, ar.StartDate, ar.EndDate, ar.Adjusted, ar.Limit, c.PolygonKey)
}

func getAggregateData(c cf.Secrets, ar models.AggregateRequest) models.PolygonAggregateResponse {
	r := parseAggregateRequest(c, ar)

	response, err := http.Get(r)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var par models.PolygonAggregateResponse
	err = json.Unmarshal(responseData, &par)

	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	return par
}

func convertAggJSONToCSV(ar models.AggregateRequest, par models.PolygonAggregateResponse) {
	folderPath := fmt.Sprintf("./data/%s/%s/", ar.Ticker, ar.StartDate)
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	destination := fmt.Sprintf("%s%s-%s-%s.csv", folderPath, ar.Ticker, ar.RangeType, ar.StartDate)
	outputFile, err := os.Create(destination)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"C", "H", "L", "N", "O", "T", "V", "Wv"}
	err = writer.Write(header)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	for _, obj := range par.Results {
		r := []string{
			fmt.Sprintf("%f", obj.C),
			fmt.Sprintf("%f", obj.H),
			fmt.Sprintf("%f", obj.L),
			fmt.Sprintf("%d", obj.N),
			fmt.Sprintf("%f", obj.O),
			fmt.Sprintf("%d", obj.T),
			fmt.Sprintf("%f", obj.V),
			fmt.Sprintf("%f", obj.Wv),
		}
		err = writer.Write(r)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	}
}

func main() {

	c := cf.GetSecrets()
	ar := models.AggregateRequest{
		Ticker:      "AAPL",
		RangeLength: 1,
		RangeType:   "minute",
		StartDate:   "2023-08-18",
		EndDate:     "2023-08-18",
		Adjusted:    true,
		Limit:       5000,
	}

	par := getAggregateData(c, ar)

	fmt.Println(par)

	convertAggJSONToCSV(ar, par)
}
