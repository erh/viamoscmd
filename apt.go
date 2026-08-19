package viamsystem

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

var AptModel = family.WithModel("apt")

func init() {
	resource.RegisterComponent(
		sensor.API,
		AptModel,
		resource.Registration[sensor.Sensor, *aptConfig]{
			Constructor: newAptSensor,
		})
}

// debian package names: lowercase letters, digits, and +-. after the first character
var debianPackageNameRegexp = regexp.MustCompile(`^[a-z0-9][a-z0-9+.\-]+$`)

type aptConfig struct {
	Packages []string
	Update   bool // run apt-get update before installing
}

func (cfg aptConfig) Validate(path string) ([]string, []string, error) {
	if len(cfg.Packages) == 0 {
		return nil, nil, fmt.Errorf("need at least one package")
	}
	for _, p := range cfg.Packages {
		if !debianPackageNameRegexp.MatchString(p) {
			return nil, nil, fmt.Errorf("invalid package name [%s]", p)
		}
	}
	return nil, nil, nil
}

func (cfg aptConfig) missingPackages() ([]string, error) {
	missing := []string{}
	for _, p := range cfg.Packages {
		installed, _, err := aptPackageStatus(p)
		if err != nil {
			return nil, err
		}
		if !installed {
			missing = append(missing, p)
		}
	}
	return missing, nil
}

func aptPackageStatus(pkg string) (bool, string, error) {
	c := exec.Command("dpkg-query", "-W", "-f", "${db:Status-Status} ${Version}", pkg)
	out, err := c.CombinedOutput()
	if err != nil {
		// dpkg-query exits non-zero for unknown packages
		if _, ok := err.(*exec.ExitError); ok {
			return false, "", nil
		}
		return false, "", fmt.Errorf("cannot check package [%s]: %w", pkg, err)
	}
	outs := strings.TrimSpace(string(out))
	status, version, _ := strings.Cut(outs, " ")
	return status == "installed", version, nil
}

func (cfg aptConfig) install(ctx context.Context, logger logging.Logger) error {
	missing, err := cfg.missingPackages()
	if err != nil {
		return err
	}
	if len(missing) == 0 {
		return nil
	}

	if cfg.Update {
		logger.Infof("running apt-get update")
		if err := runApt(ctx, "update"); err != nil {
			return err
		}
	}

	logger.Infof("installing packages %v", missing)
	args := append([]string{"install", "-y"}, missing...)
	if err := runApt(ctx, args...); err != nil {
		return err
	}
	logger.Infof("installed packages %v", missing)
	return nil
}

func runApt(ctx context.Context, args ...string) error {
	c := exec.CommandContext(ctx, "apt-get", args...)
	c.Env = append(c.Environ(), "DEBIAN_FRONTEND=noninteractive")
	out, err := c.CombinedOutput()
	if err != nil {
		outs := strings.TrimSpace(string(out))
		if len(outs) > 0 {
			return fmt.Errorf("apt-get %s failed: %s - %w", args[0], outs, err)
		}
		return fmt.Errorf("apt-get %s failed: %w", args[0], err)
	}
	return nil
}

func newAptSensor(ctx context.Context, deps resource.Dependencies, config resource.Config, logger logging.Logger) (sensor.Sensor, error) {
	newConf, err := resource.NativeConfig[*aptConfig](config)
	if err != nil {
		return nil, err
	}

	bgCtx, cancel := context.WithCancel(context.Background())
	s := &aptSensor{name: config.ResourceName(), config: newConf, logger: logger, bgCtx: bgCtx, cancel: cancel}

	s.startInstall()

	return s, nil
}

type aptSensor struct {
	resource.AlwaysRebuild

	name   resource.Name
	logger logging.Logger
	config *aptConfig
	bgCtx  context.Context
	cancel context.CancelFunc

	mu         sync.Mutex
	installing bool
	lastErr    error
}

// startInstall kicks off an install in the background so a slow apt run doesn't block startup.
func (as *aptSensor) startInstall() {
	as.mu.Lock()
	defer as.mu.Unlock()
	if as.installing {
		return
	}
	as.installing = true

	go func() {
		err := as.config.install(as.bgCtx, as.logger)
		if err != nil {
			as.logger.Errorf("error installing packages: %v", err)
		}
		as.mu.Lock()
		defer as.mu.Unlock()
		as.installing = false
		as.lastErr = err
	}()
}

func (as *aptSensor) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	allInstalled := true
	for _, p := range as.config.Packages {
		installed, version, err := aptPackageStatus(p)
		if err != nil {
			return nil, err
		}
		if installed {
			res[p] = version
		} else {
			res[p] = "missing"
			allInstalled = false
		}
	}
	res["all_installed"] = allInstalled

	as.mu.Lock()
	defer as.mu.Unlock()
	res["installing"] = as.installing
	if as.lastErr != nil {
		res["error"] = as.lastErr.Error()
	}
	return res, nil
}

func (as *aptSensor) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	if cmd["install"] == true {
		as.startInstall()
		return map[string]interface{}{"installing": true}, nil
	}
	return as.Readings(ctx, cmd)
}

func (as *aptSensor) Status(ctx context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

func (as *aptSensor) Close(ctx context.Context) error {
	as.cancel()
	return nil
}

func (as *aptSensor) Name() resource.Name {
	return as.name
}
