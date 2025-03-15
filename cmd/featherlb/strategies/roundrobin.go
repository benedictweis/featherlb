package strategies

import (
	"errors"
	"featherlb/cmd/featherlb/types"
	"sync/atomic"
)

var ErrNoBackends = errors.New("no backends available")

type RoundRobinStrategy struct {
	backends []types.Backend
	index    uint64 // Use an atomic counter
}

func NewRoundRobinStrategy() *RoundRobinStrategy {
	return &RoundRobinStrategy{
		backends: []types.Backend{},
		index:    0,
	}
}

func (r *RoundRobinStrategy) AddBackend(backend types.Backend) {
	r.backends = append(r.backends, backend)
}

func (r *RoundRobinStrategy) Next() (types.Backend, error) {
	if len(r.backends) == 0 {
		return types.Backend{}, ErrNoBackends
	}

	// Atomically increment the index and wrap around using modulo
	idx := atomic.AddUint64(&r.index, 1)
	return r.backends[(idx-1)%uint64(len(r.backends))], nil
}
