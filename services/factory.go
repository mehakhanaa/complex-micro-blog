package services

import "github.com/mehakhanaa/complex-micro-blog/stores"

type Factory struct {
	storeFactory *stores.Factory
}

func NewFactory(storeFactory *stores.Factory) *Factory {
	return &Factory{storeFactory: storeFactory}
}
