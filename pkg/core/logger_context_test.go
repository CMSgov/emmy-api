package core

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestContextHandler(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, nil)
	ctxHandler := &contextHandler{Handler: handler}
	logger := slog.New(ctxHandler)

	t.Run("with request id", func(t *testing.T) {
		buf.Reset()
		rid := "test-request-id"
		ctx := context.WithValue(context.Background(), RequestContextKey, rid)
		logger.InfoContext(ctx, "test message")

		var logRecord map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &logRecord); err != nil {
			t.Fatalf("failed to unmarshal log record: %v", err)
		}

		if logRecord["request_id"] != rid {
			t.Errorf("expected request_id %s, got %v", rid, logRecord["request_id"])
		}
		if logRecord["msg"] != "test message" {
			t.Errorf("expected msg 'test message', got %v", logRecord["msg"])
		}
	})

	t.Run("with request id and WithAttrs", func(t *testing.T) {
		buf.Reset()
		rid := "test-request-id-with-attrs"
		ctx := context.WithValue(context.Background(), RequestContextKey, rid)
		loggerWith := logger.With(slog.String("foo", "bar"))
		loggerWith.InfoContext(ctx, "test message")

		var logRecord map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &logRecord); err != nil {
			t.Fatalf("failed to unmarshal log record: %v", err)
		}

		if logRecord["request_id"] != rid {
			t.Errorf("expected request_id %s, got %v", rid, logRecord["request_id"])
		}
		if logRecord["foo"] != "bar" {
			t.Errorf("expected foo 'bar', got %v", logRecord["foo"])
		}
	})

	t.Run("without request id", func(t *testing.T) {
		buf.Reset()
		logger.InfoContext(context.Background(), "test message")

		var logRecord map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &logRecord); err != nil {
			t.Fatalf("failed to unmarshal log record: %v", err)
		}

		if _, ok := logRecord["request_id"]; ok {
			t.Error("expected no request_id in log record")
		}
	})
}
