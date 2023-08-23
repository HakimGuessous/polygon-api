package jobfetcher

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func Test_GetJobHistory(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected JobHistory
	}{
		{
			name:     "nil test",
			input:    "./job_history_test/no_file.json",
			expected: JobHistory{},
		},
		{
			name:  "get job history",
			input: "./job_history_test/test_file.json",
			expected: JobHistory{
				Jobs: []Job{
					{
						Ticker: "AAPL",
						Date:   time.Date(2023, time.August, 21, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		result := GetJobHistory(tt.input)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("For %s input %v, expected %v but got %v", tt.name, tt.input, tt.expected, result)
		}
	}
}

func Test_getNextTradingDay(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "Monday returns Tuesday",
			input:    time.Date(2023, time.August, 21, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2023, time.August, 22, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Friday returns Monday",
			input:    time.Date(2023, time.August, 18, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2023, time.August, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "Saturday returns Monday",
			input:    time.Date(2023, time.August, 19, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2023, time.August, 21, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		result := getNextTradingDay(tt.input)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("For %s input %v, expected %v but got %v", tt.name, tt.input, tt.expected, result)
		}
	}
}

func Test_WriteJobHistory(t *testing.T) {
	tests := []struct {
		name     string
		jh       JobHistory
		input    string
		expected bool
	}{
		{
			name: "write job history",
			jh: JobHistory{
				Jobs: []Job{
					{
						Ticker: "AAPL",
						Date:   time.Date(2023, time.August, 21, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			input:    "./job_history_test/test_file.json",
			expected: true,
		},
	}

	for _, tt := range tests {
		tt.jh.WriteJobHistory(tt.input)
		result := false
		if _, err := os.Stat(tt.input); err == nil {
			result = true
		}
		if result != tt.expected {
			t.Errorf("For %s input %v, expected %v but got %v", tt.name, tt.input, tt.expected, result)
		}
	}
}

func Test_GetTickerDate(t *testing.T) {
	tests := []struct {
		name     string
		jh       JobHistory
		input    string
		expected time.Time
	}{
		{
			name: "nil test",
			jh: JobHistory{
				Jobs: []Job{},
			},
			input:    "AAPL",
			expected: time.Time{},
		},
		{
			name: "get ticker max date from job history",
			jh: JobHistory{
				Jobs: []Job{
					{
						Ticker: "AAPL",
						Date:   time.Date(2023, time.August, 21, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			input:    "AAPL",
			expected: time.Date(2023, time.August, 21, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		result := tt.jh.GetTickerDate(tt.input)
		if result != tt.expected {
			t.Errorf("For %s input %v, expected %v but got %v", tt.name, tt.input, tt.expected, result)
		}
	}
}

func Test_UpdateJobHistory(t *testing.T) {
	tests := []struct {
		name     string
		jh       JobHistory
		input    []Job
		expected JobHistory
	}{
		{
			name: "nil test",
			jh: JobHistory{
				Jobs: []Job{
					{
						Ticker: "AAPL",
						Date:   time.Date(2023, time.August, 21, 0, 0, 0, 0, time.UTC),
					},
				}},
			input: []Job{},
			expected: JobHistory{
				Jobs: []Job{
					{
						Ticker: "AAPL",
						Date:   time.Date(2023, time.August, 21, 0, 0, 0, 0, time.UTC),
					},
				}},
		},
		{
			name: "Current record gets updated",
			jh: JobHistory{
				Jobs: []Job{
					{
						Ticker: "AAPL",
						Date:   time.Date(2023, time.August, 21, 0, 0, 0, 0, time.UTC),
					},
				}},
			input: []Job{
				{
					Ticker: "AAPL",
					Date:   time.Date(2023, time.August, 22, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: JobHistory{
				Jobs: []Job{
					{
						Ticker: "AAPL",
						Date:   time.Date(2023, time.August, 22, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
		{
			name: "Many records gets updated and inserted",
			jh: JobHistory{
				Jobs: []Job{
					{
						Ticker: "AAPL",
						Date:   time.Date(2023, time.August, 21, 0, 0, 0, 0, time.UTC),
					},
				}},
			input: []Job{
				{
					Ticker: "AAPL",
					Date:   time.Date(2023, time.August, 22, 0, 0, 0, 0, time.UTC),
				},
				{
					Ticker: "MSFT",
					Date:   time.Date(2023, time.August, 22, 0, 0, 0, 0, time.UTC),
				},
			},
			expected: JobHistory{
				Jobs: []Job{
					{
						Ticker: "AAPL",
						Date:   time.Date(2023, time.August, 22, 0, 0, 0, 0, time.UTC),
					},
					{
						Ticker: "MSFT",
						Date:   time.Date(2023, time.August, 22, 0, 0, 0, 0, time.UTC),
					},
				}},
		},
	}

	for _, tt := range tests {
		result := tt.jh.UpdateJobHistory(tt.input)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("For %s input %v, expected %v but got %v", tt.name, tt.input, tt.expected, result)
		}
	}
}

func Test_GetNewJobs(t *testing.T) {
	tests := []struct {
		name         string
		tickers      []string
		numberOfJobs int
		startDate    time.Time
		endDate      time.Time
		jh           JobHistory
		expected     []Job
	}{
		{
			name:         "nil AAPL returns first 2 trading days",
			tickers:      []string{"AAPL"},
			numberOfJobs: 2,
			startDate:    time.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
			endDate:      time.Date(2023, time.August, 30, 0, 0, 0, 0, time.UTC),
			jh: JobHistory{
				Jobs: []Job{},
			},
			expected: []Job{
				{
					Ticker: "AAPL",
					Date:   time.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Ticker: "AAPL",
					Date:   time.Date(2023, time.August, 2, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:         "AAPL returns next trading days",
			tickers:      []string{"AAPL"},
			numberOfJobs: 2,
			startDate:    time.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
			endDate:      time.Date(2023, time.August, 30, 0, 0, 0, 0, time.UTC),
			jh: JobHistory{
				Jobs: []Job{
					{
						Ticker: "AAPL",
						Date:   time.Date(2023, time.August, 21, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			expected: []Job{
				{
					Ticker: "AAPL",
					Date:   time.Date(2023, time.August, 22, 0, 0, 0, 0, time.UTC),
				},
				{
					Ticker: "AAPL",
					Date:   time.Date(2023, time.August, 23, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name:         "AAPL + MSFT returns last trading days",
			tickers:      []string{"AAPL", "MSFT"},
			numberOfJobs: 2,
			startDate:    time.Date(2023, time.August, 1, 0, 0, 0, 0, time.UTC),
			endDate:      time.Date(2023, time.August, 30, 0, 0, 0, 0, time.UTC),
			jh: JobHistory{
				Jobs: []Job{
					{
						Ticker: "AAPL",
						Date:   time.Date(2023, time.August, 29, 0, 0, 0, 0, time.UTC),
					},
					{
						Ticker: "MSFT",
						Date:   time.Date(2023, time.August, 29, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			expected: []Job{
				{
					Ticker: "AAPL",
					Date:   time.Date(2023, time.August, 30, 0, 0, 0, 0, time.UTC),
				},
				{
					Ticker: "MSFT",
					Date:   time.Date(2023, time.August, 30, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, tt := range tests {
		result := GetNewJobs(tt.tickers, tt.numberOfJobs, tt.startDate, tt.endDate, tt.jh)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("For %v, expected %v but got %v", tt.name, tt.expected, result)
		}
	}
}
