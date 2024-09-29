package viamoscmd

import (
	"context"
	"fmt"
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
		resource.Registration[sensor.Sensor, resource.NoNativeConfig]{
			Constructor: newCmdSensor,
		})
}

func newCmdSensor(ctx context.Context, deps resource.Dependencies, config resource.Config, logger logging.Logger) (sensor.Sensor, error) {
	s := &cmdSensor{name: config.ResourceName(), logger: logger}

	cmd, ok := config.Attributes["cmd"]
	if !ok {
		return nil, fmt.Errorf("need to specify a cmd")
	}

	cmdString, csok := cmd.(string)
	cmdArray, caok := cmd.([]string)

	if !csok && !caok {
		return nil, fmt.Errorf("cmd needs to be a string or string array, not %T", cmd)
	}

	if csok {
		cmdArray = strings.Split(cmdString, " ")
	}

	s.cmd = cmdArray[0]
	s.args = cmdArray[1:]

	argsRaw, cok := config.Attributes["args"]
	if cok {
		args, cok := argsRaw.([]string)
		if !cok {
			return nil, fmt.Errorf("args has to be a string arrray")
		}
		if len(s.args) > 0 && len(args) > 0 {
			return nil, fmt.Errorf("cannot have args in cmd and args")
		}
		s.args = args
	}

	return s, nil
}

type cmdSensor struct {
	resource.AlwaysRebuild

	name   resource.Name
	logger logging.Logger

	cmd  string
	args []string
}

func (cs *cmdSensor) run() (map[string]interface{}, error) {
	panic(1)
}

func (cs *cmdSensor) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	res, err := cs.run()
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
