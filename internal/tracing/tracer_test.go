package tracing

import (
	"context"
	"testing"
)

func TestInit_Disabled(t *testing.T) {
	shutdown, err := Init("test-service", "")
	if err != nil {
		t.Fatalf("failed to init with empty endpoint: %v", err)
	}
	defer shutdown()

	if tracer != nil {
		t.Error("expected tracer to be nil when disabled")
	}
}

func TestStartSpan_NoTracer(t *testing.T) {
	tracer = nil
	ctx := context.Background()

	newCtx, span := StartSpan(ctx, "test-span")
	if newCtx == nil {
		t.Error("expected context, got nil")
	}
	if span == nil {
		t.Error("expected span, got nil")
	}
}

func TestGetTracer(t *testing.T) {
	tracer = nil
	tr := GetTracer()
	if tr != nil {
		t.Error("expected nil tracer")
	}
}
