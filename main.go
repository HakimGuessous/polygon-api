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

	jh := jobfetcher.GetJobHistory("./job_history/job_history.json")
	jobs := jobfetcher.GetNewJobs(c.Tickers, c.NumberOfJobs, c.JobsStartDate, c.JobsEndDate, jh)

	for _, job := range jobs {
		ar := polygon.ConvertJobToAggregateRequest(job, c.JobRangeLength, c.JobRangeType, c.JobAdjusted, c.JobLimit)
		par := polygon.GetAggregateData(s.PolygonKey, ar)

		path := fmt.Sprintf("./data/%s/%s/%s-%s-%s.csv", ar.Ticker, ar.StartDate, ar.Ticker, ar.RangeType, ar.StartDate)
		polygon.ConvertAggJSONToCSV(par, path)
	}
	jh = jh.UpdateJobHistory(jobs)
	jh.WriteJobHistory("./job_history/job_history.json")
}
