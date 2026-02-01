package config

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// Watcher watches configuration file for changes
type Watcher struct {
	configPath string
	watcher    *fsnotify.Watcher
	onChange   func(*Config) error
	onError    func(error)
	mu         sync.Mutex
	stopCh     chan struct{}
}

// NewWatcher creates a new configuration file watcher
func NewWatcher(configPath string, onChange func(*Config) error, onError func(error)) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	if err := w.Add(configPath); err != nil {
		_ = w.Close()
		return nil, fmt.Errorf("failed to watch config file: %w", err)
	}

	return &Watcher{
		configPath: configPath,
		watcher:    w,
		onChange:   onChange,
		onError:    onError,
		stopCh:     make(chan struct{}),
	}, nil
}

// Start begins watching for configuration changes
func (w *Watcher) Start() {
	go w.watch()
}

// Stop stops watching for configuration changes
func (w *Watcher) Stop() {
	close(w.stopCh)
	_ = w.watcher.Close()
}

func (w *Watcher) watch() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				w.handleChange()
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			if w.onError != nil {
				w.onError(fmt.Errorf("watcher error: %w", err))
			}
		case <-w.stopCh:
			return
		}
	}
}

func (w *Watcher) handleChange() {
	w.mu.Lock()
	defer w.mu.Unlock()

	cfg, err := Load(w.configPath)
	if err != nil {
		if w.onError != nil {
			w.onError(fmt.Errorf("failed to reload config: %w", err))
		}
		return
	}

	if err := w.onChange(cfg); err != nil {
		if w.onError != nil {
			w.onError(fmt.Errorf("failed to apply config changes: %w", err))
		}
	}
}
