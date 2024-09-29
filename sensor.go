package viamoscmd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

var family = resource.ModelNamespace("erh").WithFamily("viamoscmd")

var Model = family.WithModel("cmdsensor")

func init() {
	resource.RegisterComponent(
		sensor.API,
		Model,
		resource.Registration[sensor.Sensor, *cmdSensorConfig]{
			Constructor: newCmdSensor,
		})
}

type cmdSensorConfig struct {
	Cmd string
	Args []string
}

func (cfg cmdSensorConfig) Validate(path string) ([]string, error) {
	if cfg.Cmd == "" {
		return nil, fmt.Errorf("need cmd")
	}
	if strings.Index(cfg.Cmd, " ") >= 0 {
		return nil, fmt.Errorf("cmd cannot have spaces")
	}
	return nil, nil
}

func (cfg cmdSensorConfig) run() (map[string]interface{}, error) {
	c := exec.Command(cfg.Cmd, cfg.Args...)
	out, err := c.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"out" : string(out)}, nil
}


func newCmdSensor(ctx context.Context, deps resource.Dependencies, config resource.Config, logger logging.Logger) (sensor.Sensor, error) {
	newConf, err := resource.NativeConfig[*cmdSensorConfig](config)
	if err != nil {
		return nil, err
	}

	s := &cmdSensor{name: config.ResourceName(), config: newConf, logger: logger}

	return s, nil
}

type cmdSensor struct {
	resource.AlwaysRebuild

	name   resource.Name
	logger logging.Logger
	config *cmdSensorConfig
}

func (cs *cmdSensor) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	res, err := cs.config.run()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (cs *cmdSensor) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (cs *cmdSensor) Close(ctx context.Context) error {
	return nil
}

func (cs *cmdSensor) Name() resource.Name {
	return cs.name
}
