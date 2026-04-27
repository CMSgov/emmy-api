package reporting

import (
	"time"
)

// ReportData represents the data structure received from SQS
type ReportData struct {
	Timestamp  time.Time `json:"timestamp"`
	Endpoint   string    `json:"endpoint"`
	DataSource string    `json:"data_source"` //nolint:tagliatelle // Reporting payload contract uses snake_case.
	ClientID   string    `json:"client_id"`   //nolint:tagliatelle // Reporting payload contract uses snake_case.
	StatusCode int       `json:"status_code"` //nolint:tagliatelle // Reporting payload contract uses snake_case.
	Success    bool      `json:"success"`
}
