package strategies

import (
	"errors"
	"featherlb/cmd/featherlb/types"
	"sync"
)

var ErrNoBackends = errors.New("no backends available")

type RoundRobinStrategy struct {
	backends []types.Backend
	index    int
	mu       sync.Mutex // Mutex to ensure thread safety
}

func NewRoundRobinStrategy() *RoundRobinStrategy {
	return &RoundRobinStrategy{
		backends: []types.Backend{},
		index:    0,
	}
}

func (r *RoundRobinStrategy) AddBackend(backend types.Backend) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.backends = append(r.backends, backend)
}

func (r *RoundRobinStrategy) Next() (types.Backend, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.backends) == 0 {
		return types.Backend{}, ErrNoBackends
	}

	if r.index >= len(r.backends) {
		r.index = 0
	}

	backend := r.backends[r.index]
	r.index = (r.index + 1) % len(r.backends)

	return backend, nil
}
