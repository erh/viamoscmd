package viamsystem

import (
	"os/exec"
	"testing"

	"go.viam.com/test"
)

func TestAptConfigValidate(t *testing.T) {
	cfg := aptConfig{}
	_, _, err := cfg.Validate("")
	test.That(t, err, test.ShouldNotBeNil)

	cfg.Packages = []string{"curl"}
	_, _, err = cfg.Validate("")
	test.That(t, err, test.ShouldBeNil)

	cfg.Packages = []string{"curl", "g++", "libc6-dev", "python3.11"}
	_, _, err = cfg.Validate("")
	test.That(t, err, test.ShouldBeNil)

	cfg.Packages = []string{"curl; rm -rf /"}
	_, _, err = cfg.Validate("")
	test.That(t, err, test.ShouldNotBeNil)

	cfg.Packages = []string{"Curl"}
	_, _, err = cfg.Validate("")
	test.That(t, err, test.ShouldNotBeNil)

	cfg.Packages = []string{""}
	_, _, err = cfg.Validate("")
	test.That(t, err, test.ShouldNotBeNil)
}

func TestAptPackageStatus(t *testing.T) {
	if _, err := exec.LookPath("dpkg-query"); err != nil {
		t.Skip("dpkg-query not available")
	}

	installed, _, err := aptPackageStatus("surely-not-a-real-package-xyz")
	test.That(t, err, test.ShouldBeNil)
	test.That(t, installed, test.ShouldBeFalse)

	installed, version, err := aptPackageStatus("dpkg")
	test.That(t, err, test.ShouldBeNil)
	test.That(t, installed, test.ShouldBeTrue)
	test.That(t, version, test.ShouldNotBeEmpty)
}
