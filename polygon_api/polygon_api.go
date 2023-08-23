package polygonapi

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	jobfetcher "polygon/job_fetcher"
)

func parseAggregateRequest(polygonKey string, ar AggregateRequest) string {
	return fmt.Sprintf("https://api.polygon.io/v2/aggs/ticker/%s/range/%d/%s/%s/%s?adjusted=%t&sort=asc&limit=%d&apiKey=%s",
		ar.Ticker, ar.RangeLength, ar.RangeType, ar.StartDate, ar.EndDate, ar.Adjusted, ar.Limit, polygonKey)
}

func ConvertJobToAggregateRequest(job jobfetcher.Job, rangeLength int, rangeType string, adjusted bool, limit int) AggregateRequest {
	ar := AggregateRequest{
		Ticker:      job.Ticker,
		RangeLength: rangeLength,
		RangeType:   rangeType,
		StartDate:   job.Date.Format("2006-01-02"),
		EndDate:     job.Date.Format("2006-01-02"),
		Adjusted:    adjusted,
		Limit:       limit,
	}

	return ar
}

func GetAggregateData(polygonKey string, ar AggregateRequest) PolygonAggregateResponse {
	r := parseAggregateRequest(polygonKey, ar)

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

	var par PolygonAggregateResponse
	err = json.Unmarshal(responseData, &par)

	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	return par
}

func ConvertAggJSONToCSV(par PolygonAggregateResponse, path string) {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	outputFile, err := os.Create(path)
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

type AggregateRequest struct {
	Ticker      string
	RangeLength int
	RangeType   string
	StartDate   string
	EndDate     string
	Adjusted    bool
	Limit       int
}

type PolygonAggregateResponse struct {
	Adjusted     bool
	Next_url     string
	QueryCount   int
	Request_id   string
	Results      []AggregateResult
	ResultsCount int
	Status       string
	Ticker       string
	Count        int
}

type AggregateResult struct {
	C  float32
	H  float32
	L  float32
	N  int64
	O  float32
	T  int64
	V  float32
	Wv float32
}
