package shutdown

import (
	"errors"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestNewHandler(t *testing.T) {
	handler := NewHandler(5 * time.Second)
	if handler == nil {
		t.Fatal("expected handler, got nil")
	}
	defer handler.Stop()

	if handler.timeout != 5*time.Second {
		t.Errorf("expected timeout 5s, got %v", handler.timeout)
	}
}

func TestRegister(t *testing.T) {
	handler := NewHandler(5 * time.Second)
	defer handler.Stop()

	called := false
	handler.Register(func() error {
		called = true
		return nil
	})

	if err := handler.Shutdown(); err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}

	if !called {
		t.Error("expected handler to be called")
	}
}

func TestShutdown_Success(t *testing.T) {
	handler := NewHandler(5 * time.Second)
	defer handler.Stop()

	count := 0
	handler.Register(func() error {
		count++
		return nil
	})
	handler.Register(func() error {
		count++
		return nil
	})

	if err := handler.Shutdown(); err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 handlers called, got %d", count)
	}
}

func TestShutdown_WithError(t *testing.T) {
	handler := NewHandler(5 * time.Second)
	defer handler.Stop()

	testErr := errors.New("test error")
	handler.Register(func() error {
		return testErr
	})

	if err := handler.Shutdown(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestShutdown_Timeout(t *testing.T) {
	handler := NewHandler(100 * time.Millisecond)
	defer handler.Stop()

	handler.Register(func() error {
		time.Sleep(500 * time.Millisecond)
		return nil
	})

	err := handler.Shutdown()
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}

func TestWait_Signal(t *testing.T) {
	handler := NewHandler(5 * time.Second)
	defer handler.Stop()

	go func() {
		time.Sleep(100 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(syscall.SIGTERM)
	}()

	select {
	case sig := <-handler.Wait():
		if sig != syscall.SIGTERM {
			t.Errorf("expected SIGTERM, got %v", sig)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for signal")
	}
}

func TestShutdown_MultipleHandlers(t *testing.T) {
	handler := NewHandler(5 * time.Second)
	defer handler.Stop()

	order := []int{}
	handler.Register(func() error {
		order = append(order, 1)
		return nil
	})
	handler.Register(func() error {
		order = append(order, 2)
		return nil
	})
	handler.Register(func() error {
		order = append(order, 3)
		return nil
	})

	if err := handler.Shutdown(); err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}

	if len(order) != 3 {
		t.Errorf("expected 3 handlers, got %d", len(order))
	}

	for i, v := range order {
		if v != i+1 {
			t.Errorf("expected order[%d] = %d, got %d", i, i+1, v)
		}
	}
}
