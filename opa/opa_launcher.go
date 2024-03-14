package opa

import (
	"fmt"
	"os"
	"os/exec"
	"raygun/log"
	"time"
)

type OpaConfig struct {
	OpaPort    uint16
	OpaPath    string
	BundlePath string
	LogPath    string
}

func (oc OpaConfig) String() string {
	return fmt.Sprintf("exec: %s, bundle: %s, logs: %s", oc.OpaPath, oc.BundlePath, oc.LogPath)
}

type OpaRunner struct {
	Config  OpaConfig
	Process *os.Process
}

func NewOpaRunner(config OpaConfig) OpaRunner {

	log.Debug("Building new OpaRunner with config: %v", config)
	tmp := OpaRunner{Config: config}

	return tmp

}

func (opa *OpaRunner) Start() error {

	commandToRun := opa.Config.OpaPath

	absolute_path, err := exec.LookPath(commandToRun)

	if err != nil {
		log.Error("Unable to find %s on the path", commandToRun)
	}

	log.Debug("OpaRunner.Start() - commandToRun: %s - absolute_path: %s", commandToRun, absolute_path)

	args := []string{commandToRun, "run", "--server", "-b", opa.Config.BundlePath}

	log.Debug("OpaRunner.Start() - arg string: %v", args)

	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	process, err := os.StartProcess(absolute_path, args, procAttr)

	if err != nil {
		log.Error("Unable to start OPA: %s", err.Error())
		return err
	}

	log.Debug("Started %s with process id: %d", commandToRun, process.Pid)

	opa.Process = process

	time.Sleep(1 * time.Second)

	log.Debug("Hopefully OPA is now up")

	return nil
}

func (opa *OpaRunner) Stop() error {

	if opa.Process == nil {
		return fmt.Errorf("OpaRunner:Stop - no process found, can't stop, won't stop")
	}

	log.Debug("OpaRunner:Stop() - stoppping process: %d", opa.Process.Pid)

	opa.Process.Kill()

	return nil
}
