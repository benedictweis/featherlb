package strategies

import (
	"errors"
	"featherlb/internal/pkg/types"
	"sync/atomic"
)

var ErrNoEndpoints = errors.New("no endpoints available")

type RoundRobinStrategy struct {
	endpoints []types.Endpoint
	index     uint64 // Use an atomic counter
}

func NewRoundRobinStrategy() *RoundRobinStrategy {
	return &RoundRobinStrategy{
		endpoints: []types.Endpoint{},
		index:     0,
	}
}

func (r *RoundRobinStrategy) AddEndpoint(endpoint types.Endpoint) {
	r.endpoints = append(r.endpoints, endpoint)
}

func (r *RoundRobinStrategy) Next() (types.Endpoint, error) {
	if len(r.endpoints) == 0 {
		return types.Endpoint{}, ErrNoEndpoints
	}

	// Atomically increment the index and wrap around using modulo
	idx := atomic.AddUint64(&r.index, 1)
	return r.endpoints[(idx-1)%uint64(len(r.endpoints))], nil
}
