package viamoscmd

import (
	"testing"

	"go.viam.com/test"
)

func TestCmdSensorErrors(t *testing.T) {
	cfg := cmdSensorConfig{}
	_, err := cfg.Validate("")
	test.That(t, err, test.ShouldNotBeNil)

	cfg.Cmd = "echo 1"
	_, err = cfg.Validate("")
	test.That(t, err, test.ShouldNotBeNil)

	cfg.Cmd = "echo"
	_, err = cfg.Validate("")
	test.That(t, err, test.ShouldBeNil)

	cfg.Cmd = "echo"
	_, err = cfg.Validate("")
	test.That(t, err, test.ShouldBeNil)
	res, err := cfg.run()
	test.That(t, err, test.ShouldBeNil)
	test.That(t, map[string]interface{}{"out" : "\n"}, test.ShouldResemble, res)

	
	cfg.Cmd = "echo"
	cfg.Args = []string{"1"}
	_, err = cfg.Validate("")
	test.That(t, err, test.ShouldBeNil)
	res, err = cfg.run()
	test.That(t, err, test.ShouldBeNil)
	test.That(t, map[string]interface{}{"out" : "1\n"}, test.ShouldResemble, res)
}
