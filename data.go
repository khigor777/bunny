package bunny

import "context"

type Set struct {
	Key   string
	Value []byte
	Ctx   context.Context
}
