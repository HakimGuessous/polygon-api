package main

import (
	"fmt"
	"time"

	cf "polygon/config"
	jobfetcher "polygon/job_fetcher"
	polygon "polygon/polygon_api"
)

func getData() []jobfetcher.Job {
	s := cf.GetSecrets()
	c := cf.GetConfig()

	jh := jobfetcher.GetJobHistory("./job_history/job_history.json")
	jobs := jobfetcher.GetNewJobs(c.Tickers, c.NumberOfJobs, c.JobsStartDate, c.JobsEndDate, jh)

	for _, job := range jobs {
		ar := polygon.ConvertJobToAggregateRequest(job, c.JobRangeLength, c.JobRangeType, c.JobAdjusted, c.JobLimit)
		par := polygon.GetAggregateData(s.PolygonKey, ar)

		path := fmt.Sprintf("./data/%s/%s/%s-%s-%s.csv", ar.Ticker, ar.StartDate, ar.Ticker, ar.RangeType, ar.StartDate)
		polygon.ConvertAggJSONToCSV(par, path)
		print("Written file: %s - %s", ar.Ticker, ar.StartDate)
	}
	jh = jh.UpdateJobHistory(jobs)
	jh.WriteJobHistory("./job_history/job_history.json")

	return jobs
}

func main() {
	ticker := time.NewTicker(1 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				jobs := getData()
				if len(jobs) == 0 {
					close(quit)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
