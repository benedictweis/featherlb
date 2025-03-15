package strategies

import "featherlb/cmd/featherlb/types"

type Strategy interface {
	AddBackend(backend types.Backend)
	Next() (backend types.Backend, err error)
}
