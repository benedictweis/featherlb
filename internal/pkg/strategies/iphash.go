package strategies

import (
	"featherlb/internal/pkg/types"
	"hash/fnv"
	"net"
)

type IPHashStrategy struct {
	endpoints []types.Endpoint
}

func NewIPHashStrategy() *IPHashStrategy {
	return &IPHashStrategy{
		endpoints: []types.Endpoint{},
	}
}

func (i *IPHashStrategy) AddEndpoint(endpoint types.Endpoint) {
	i.endpoints = append(i.endpoints, endpoint)
}

func (i *IPHashStrategy) Next(ip net.TCPAddr) (types.Endpoint, error) {
	if len(i.endpoints) == 0 {
		return types.Endpoint{}, ErrNoEndpoints
	}

	// Hash the IP address to get an index
	hash := hashIP(ip)
	idx := hash % uint64(len(i.endpoints))
	return i.endpoints[idx], nil
}

func hashIP(ip net.TCPAddr) uint64 {
	hasher := fnv.New64a()
	_, err := hasher.Write(ip.IP)
	if err != nil {
		// Handle the unlikely case of a write error
		return 0
	}
	return hasher.Sum64()
}
