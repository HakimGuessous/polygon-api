package models

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
