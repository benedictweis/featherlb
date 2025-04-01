package strategies

import "featherlb/internal/pkg/types"

type Strategy interface {
	AddEndpoint(endpoint types.Endpoint)
	Next() (endpoint types.Endpoint, err error)
}
