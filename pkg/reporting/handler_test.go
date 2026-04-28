package reporting

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLambdaHandler_HandleRequest(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	data := ReportData{
		Endpoint:   "/test",
		DataSource: "test-source",
		ClientID:   "test-client",
		Success:    true,
		Timestamp:  now,
		StatusCode: 200,
	}
	body, err := json.Marshal(data)
	require.NoError(t, err)

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId: "1",
				Body:      string(body),
			},
		},
	}

	handler := NewLambdaHandler(nil, nil)
	err = handler.HandleRequest(context.Background(), sqsEvent)

	assert.NoError(t, err)
}

func TestLambdaHandler_HandleRequest_InvalidJSON(t *testing.T) {
	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId: "2",
				Body:      "invalid json",
			},
		},
	}

	handler := NewLambdaHandler(nil, nil)
	err := handler.HandleRequest(context.Background(), sqsEvent)

	assert.NoError(t, err) // Should log error and continue, not return error
}
