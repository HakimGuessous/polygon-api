package main

import (
	"fmt"

	cf "polygon/config"
	jobfetcher "polygon/job_fetcher"
	polygon "polygon/polygon_api"
)

func main() {
	s := cf.GetSecrets()
	c := cf.GetConfig()

	jh := jobfetcher.GetJobHistory("./job_history_test/job_history.json")
	jobs := jobfetcher.GetNewJobs(c.Tickers, c.NumberOfJobs, c.JobsStartDate, c.JobsEndDate, jh)

	for _, job := range jobs {
		ar := polygon.ConvertJobToAggregateRequest(job, c.jobRangeLength, c.jobRangeType, c.jobAdjusted, c.jobLimit)
	}

	ar := polygon.AggregateRequest{
		Ticker:      "AAPL",
		RangeLength: 1,
		RangeType:   "minute",
		StartDate:   "2023-08-18",
		EndDate:     "2023-08-18",
		Adjusted:    true,
		Limit:       5000,
	}

	par := polygon.GetAggregateData(s.PolygonKey, ar)

	path := fmt.Sprintf("./data/%s/%s/%s-%s-%s.csv", ar.Ticker, ar.StartDate, ar.Ticker, ar.RangeType, ar.StartDate)
	polygon.ConvertAggJSONToCSV(par, path)
}
