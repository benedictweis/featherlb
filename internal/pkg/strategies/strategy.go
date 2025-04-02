package strategies

import (
	"errors"
	"featherlb/internal/pkg/types"
	"net"
)

var ErrNoEndpoints = errors.New("no endpoints available")

type Strategy interface {
	AddEndpoint(endpoint types.Endpoint)
	Next(ip net.TCPAddr) (endpoint types.Endpoint, err error)
}

func MatchStrategy(strategy types.KnownStrategy) Strategy {
	switch strategy {
	case types.StrategyRoundRobin:
		return NewRoundRobinStrategy()
	case types.StrategyRandom:
		return NewRandomStrategy()
	case types.StrategyIPHash:
		return NewIPHashStrategy()
	default:
		panic("unknown strategy")
	}
}
