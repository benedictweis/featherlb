package strategies

import (
	"featherlb/internal/pkg/types"
	"math/rand"
	"net"
)

type RandomStrategy struct {
	endpoints []types.Endpoint
}

func NewRandomStrategy() *RandomStrategy {
	return &RandomStrategy{
		endpoints: []types.Endpoint{},
	}
}

func (r *RandomStrategy) AddEndpoint(endpoint types.Endpoint) {
	r.endpoints = append(r.endpoints, endpoint)
}

func (r *RandomStrategy) Next(_ net.TCPAddr) (types.Endpoint, error) {
	if len(r.endpoints) == 0 {
		return types.Endpoint{}, ErrNoEndpoints
	}

	// Randomly select an endpoint
	idx := rand.Intn(len(r.endpoints))
	return r.endpoints[idx], nil
}
