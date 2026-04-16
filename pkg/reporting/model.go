package reporting

import (
	"time"
)

// ReportData represents the data structure received from SQS
type ReportData struct {
	Timestamp  time.Time `json:"timestamp"`
	Endpoint   string    `json:"endpoint"`
	DataSource string    `json:"data_source"`
	ClientID   string    `json:"client_id"`
	StatusCode int       `json:"status_code"`
	Success    bool      `json:"success"`
}
