package syncer

import (
	"context"
	"sync"
	"time"

	"github.com/ohauer/docker-secrets/internal/config"
)

// Scheduler manages periodic secret synchronization
type Scheduler struct {
	syncer  *SecretSyncer
	jobs    map[string]*job
	mu      sync.RWMutex
	stopCh  chan struct{}
	results chan SyncResult
}

type job struct {
	secret   config.Secret
	ticker   *time.Ticker
	stopCh   chan struct{}
	lastSync time.Time
}

// NewScheduler creates a new scheduler
func NewScheduler(syncer *SecretSyncer) *Scheduler {
	return &Scheduler{
		syncer:  syncer,
		jobs:    make(map[string]*job),
		stopCh:  make(chan struct{}),
		results: make(chan SyncResult, 100),
	}
}

// AddSecret adds a secret to the scheduler
func (s *Scheduler) AddSecret(cfg *config.Config, secret config.Secret) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if existing, ok := s.jobs[secret.Name]; ok {
		existing.ticker.Stop()
		close(existing.stopCh)
	}

	j := &job{
		secret: secret,
		ticker: time.NewTicker(secret.RefreshInterval),
		stopCh: make(chan struct{}),
	}

	s.jobs[secret.Name] = j

	go s.runJob(cfg, j)
}

// RemoveSecret removes a secret from the scheduler
func (s *Scheduler) RemoveSecret(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if j, ok := s.jobs[name]; ok {
		j.ticker.Stop()
		close(j.stopCh)
		delete(s.jobs, name)
	}
}

// Stop stops all scheduled jobs
func (s *Scheduler) Stop() {
	close(s.stopCh)

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, j := range s.jobs {
		j.ticker.Stop()
		close(j.stopCh)
	}
	s.jobs = make(map[string]*job)
}

// Results returns the results channel
func (s *Scheduler) Results() <-chan SyncResult {
	return s.results
}

func (s *Scheduler) runJob(cfg *config.Config, j *job) {
	ctx := context.Background()

	s.syncAndReport(ctx, cfg, j)

	for {
		select {
		case <-j.ticker.C:
			s.syncAndReport(ctx, cfg, j)
		case <-j.stopCh:
			return
		case <-s.stopCh:
			return
		}
	}
}

func (s *Scheduler) syncAndReport(ctx context.Context, cfg *config.Config, j *job) {
	err := s.syncer.SyncSecret(ctx, cfg, j.secret)

	result := SyncResult{
		SecretName: j.secret.Name,
		Success:    err == nil,
		Error:      err,
		Timestamp:  time.Now(),
	}

	if err == nil {
		j.lastSync = result.Timestamp
	}

	select {
	case s.results <- result:
	default:
	}
}

// GetLastSyncTime returns the last successful sync time for a secret
func (s *Scheduler) GetLastSyncTime(name string) (time.Time, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if j, ok := s.jobs[name]; ok {
		return j.lastSync, true
	}
	return time.Time{}, false
}
