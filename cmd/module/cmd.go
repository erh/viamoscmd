package main

import (
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"

	"github.com/erh/viamoscmd"
)

func main() {
	module.ModularMain(
		resource.APIModel{sensor.API, viamoscmd.Model},
	)
}
