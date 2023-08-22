package jobfetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func GetJobHistory(path string) JobHistory {
	var jh JobHistory

	if _, err := os.Stat(path); err == nil {
		// path/to/whatever exists
		jobHistoryFile, err := os.Open(path)

		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		defer jobHistoryFile.Close()

		jobHistoryJson, _ := io.ReadAll(jobHistoryFile)

		err = json.Unmarshal(jobHistoryJson, &jh)

		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
	}

	return jh
}

func (jh JobHistory) WriteJobHistory(path string) {
	jobHistoryJson, err := json.Marshal(jh)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = os.WriteFile(path, jobHistoryJson, 0644)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func (jh JobHistory) GetTickerMaxDate(t string) time.Time {
	var maxDate time.Time
	for _, jh := range jh.Jobs {
		if jh.Ticker == t {
			maxDate = jh.MaxDate
		}
	}
	return maxDate
}

func getNextTradingDay(t time.Time) time.Time {
	if int(t.Weekday()) == 5 || int(t.Weekday()) == 6 {
		return getNextTradingDay(t.AddDate(0, 0, 1))
	} else {
		return t.AddDate(0, 0, 1)
	}
}

func GetNextJob(t string, maxDate time.Time) Job {
	job := Job{
		Ticker:  t,
		MaxDate: getNextTradingDay(maxDate),
	}
	return job
}

func GetNewJobs(tickers []string, numberOfJobs int, startDate time.Time, endDate time.Time, jh JobHistory) []Job {
	var jobs []Job

	i := 0
	for _, t := range tickers {
		maxDate := jh.GetTickerMaxDate(t)

		for i < numberOfJobs {
			if maxDate.Equal(time.Time{}) {
				j := Job{
					Ticker:  t,
					MaxDate: startDate,
				}
				jobs = append(jobs, j)
				maxDate = startDate
				i += 1
			} else if maxDate.Before(endDate) {
				j := GetNextJob(t, maxDate)
				jobs = append(jobs, j)
				maxDate = j.MaxDate
				i += 1
			} else {
				break
			}
		}
		if i == numberOfJobs {
			break
		}
	}
	return jobs
}

type JobHistory struct {
	Jobs []Job
}

type Job struct {
	Ticker  string
	MaxDate time.Time
}
