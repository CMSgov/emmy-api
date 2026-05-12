package reporting

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeEventStore struct {
	lookupAgencyName string
	lookupErr        error
	insertData       *ReportData
	insertAgencyName *string
	insertErr        error
}

func (s *fakeEventStore) LookupAgencyName(context.Context, string) (string, error) {
	return s.lookupAgencyName, s.lookupErr
}

func (s *fakeEventStore) InsertAPIEvent(_ context.Context, data *ReportData, agencyName *string) error {
	copy := *data
	s.insertData = &copy
	s.insertAgencyName = agencyName
	return s.insertErr
}

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
		Records: []events.SQSMessage{{
			MessageId: "1",
			Body:      string(body),
		}},
	}

	store := &fakeEventStore{lookupAgencyName: "Colorado"}
	handler := &LambdaHandler{logger: slog.Default(), store: store}
	err = handler.HandleRequest(context.Background(), sqsEvent)

	assert.NoError(t, err)
	require.NotNil(t, store.insertData)
	assert.Equal(t, data, *store.insertData)
	require.NotNil(t, store.insertAgencyName)
	assert.Equal(t, "Colorado", *store.insertAgencyName)
}

func TestLambdaHandler_HandleRequest_MissingAgencyMapping(t *testing.T) {
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
		Records: []events.SQSMessage{{
			MessageId: "1",
			Body:      string(body),
		}},
	}

	store := &fakeEventStore{lookupErr: sql.ErrNoRows}
	handler := &LambdaHandler{logger: slog.Default(), store: store}
	err = handler.HandleRequest(context.Background(), sqsEvent)

	assert.NoError(t, err)
	require.NotNil(t, store.insertData)
	assert.Nil(t, store.insertAgencyName)
}

func TestLambdaHandler_HandleRequest_LookupErrorStillInserts(t *testing.T) {
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
		Records: []events.SQSMessage{{
			MessageId: "1",
			Body:      string(body),
		}},
	}

	store := &fakeEventStore{lookupErr: errors.New("lookup failed")}
	handler := &LambdaHandler{logger: slog.Default(), store: store}
	err = handler.HandleRequest(context.Background(), sqsEvent)

	assert.NoError(t, err)
	require.NotNil(t, store.insertData)
	assert.Nil(t, store.insertAgencyName)
}

func TestLambdaHandler_HandleRequest_InvalidJSON(t *testing.T) {
	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{{
			MessageId: "2",
			Body:      "invalid json",
		}},
	}

	handler := NewLambdaHandler(nil, nil)
	err := handler.HandleRequest(context.Background(), sqsEvent)

	assert.NoError(t, err)
}
