package middlewares

import "github.com/mehakhanaa/complex-micro-blog/stores"

type Factory struct {
	store *stores.Factory
}

func NewFactory(store *stores.Factory) *Factory {
	return &Factory{store: store}
}
