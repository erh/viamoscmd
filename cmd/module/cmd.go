package main

import (
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"

	viamsystem "github.com/erh/viam-system"
)

func main() {
	module.ModularMain(
		resource.APIModel{API: sensor.API, Model: viamsystem.Model},
	)
}
