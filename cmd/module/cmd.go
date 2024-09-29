package main

import (
	"context"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"

	"github.com/erh/viamoscmd"
)

func main() {
	err := realMain()
	if err != nil {
		panic(err)
	}
}
func realMain() error {

	ctx := context.Background()
	logger := logging.NewLogger("main")

	myMod, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}

	err = myMod.AddModelFromRegistry(ctx, sensor.API, viamoscmd.Model)
	if err != nil {
		return err
	}

	err = myMod.Start(ctx)
	defer myMod.Close(ctx)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}
