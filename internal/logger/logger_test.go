package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInit_ValidLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}
	for _, level := range levels {
		if err := Init(level); err != nil {
			t.Errorf("Init(%s) failed: %v", level, err)
		}
	}
}

func TestInit_DefaultLevel(t *testing.T) {
	if err := Init("invalid"); err != nil {
		t.Errorf("Init with invalid level should not fail: %v", err)
	}
}

func TestLogger_JSONOutput(t *testing.T) {
	var buf bytes.Buffer

	encoderConfig := zap.NewProductionEncoderConfig()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(&buf),
		zapcore.InfoLevel,
	)
	globalLogger = zap.New(core)

	Info("test message", zap.String("key", "value"))

	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("failed to parse JSON log: %v", err)
	}

	if logEntry["msg"] != "test message" {
		t.Errorf("expected msg 'test message', got: %v", logEntry["msg"])
	}

	if logEntry["key"] != "value" {
		t.Errorf("expected key 'value', got: %v", logEntry["key"])
	}

	if logEntry["level"] != "info" {
		t.Errorf("expected level 'info', got: %v", logEntry["level"])
	}
}

func TestLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer

	encoderConfig := zap.NewProductionEncoderConfig()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(&buf),
		zapcore.InfoLevel,
	)
	globalLogger = zap.New(core)

	childLogger := With(zap.String("component", "test"))
	childLogger.Info("child message")

	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("failed to parse JSON log: %v", err)
	}

	if logEntry["component"] != "test" {
		t.Errorf("expected component 'test', got: %v", logEntry["component"])
	}
}

func TestLogger_DifferentLevels(t *testing.T) {
	var buf bytes.Buffer

	encoderConfig := zap.NewProductionEncoderConfig()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(&buf),
		zapcore.DebugLevel,
	)
	globalLogger = zap.New(core)

	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")

	logs := bytes.Split(buf.Bytes(), []byte("\n"))
	if len(logs) < 4 {
		t.Fatalf("expected at least 4 log entries, got: %d", len(logs))
	}

	levels := []string{"debug", "info", "warn", "error"}
	for i, level := range levels {
		var logEntry map[string]interface{}
		if err := json.Unmarshal(logs[i], &logEntry); err != nil {
			t.Fatalf("failed to parse JSON log %d: %v", i, err)
		}
		if logEntry["level"] != level {
			t.Errorf("expected level '%s', got: %v", level, logEntry["level"])
		}
	}
}
