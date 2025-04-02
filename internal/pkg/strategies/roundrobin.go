package strategies

import (
	"featherlb/internal/pkg/types"
	"net"
	"sync/atomic"
)

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

func (r *RoundRobinStrategy) Next(_ net.TCPAddr) (types.Endpoint, error) {
	if len(r.endpoints) == 0 {
		return types.Endpoint{}, ErrNoEndpoints
	}

	// Atomically increment the index and wrap around using modulo
	idx := atomic.AddUint64(&r.index, 1)
	return r.endpoints[(idx-1)%uint64(len(r.endpoints))], nil
}
