package viamoscmd

import (
	"context"
	"testing"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/utils"
	"go.viam.com/test"
)

func TestCmdSensorErrors(t *testing.T) {
	ctx := context.Background()
	deps := resource.Dependencies{}
	logger := logging.NewTestLogger(t)

	_, err := newCmdSensor(ctx, deps, resource.Config{Attributes: utils.AttributeMap{}}, logger)
	test.That(t, err, test.ShouldNotBeNil)

	_, err = newCmdSensor(ctx, deps, resource.Config{Attributes: utils.AttributeMap{"cmd": 1}}, logger)
	test.That(t, err, test.ShouldNotBeNil)

	_, err = newCmdSensor(ctx, deps, resource.Config{Attributes: utils.AttributeMap{"cmd": []interface{}{"echo", 1}}}, logger)
	test.That(t, err, test.ShouldBeNil)

}
