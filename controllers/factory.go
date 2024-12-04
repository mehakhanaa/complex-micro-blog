package controllers

import (
	"github.com/mehakhanaa/complex-micro-blog/services"
)

type Factory struct {
	serviceFactory *services.Factory
}

func NewFactory(serviceFactory *services.Factory) *Factory {
	return &Factory{serviceFactory: serviceFactory}
}
