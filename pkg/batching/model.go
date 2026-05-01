package batching

import (
	"time"
)

// BatchData represents the data structure received from SQS
type BatchData struct {
	Timestamp time.Time `json:"timestamp"`
	BatchID   string    `json:"batchId"`
}
