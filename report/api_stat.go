package report

// ApiStat the statistical data of api access info
type ApiStat struct {
	// AccessCount all the api access count
	AccessCount uint64 `json:"access_count"`
	// VisitorStat the statistical data of visitors
	VisitorStat map[string]uint64 `json:"visitor_stat"`
}
