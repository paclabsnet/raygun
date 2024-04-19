/*
Copyright © 2024 PACLabs
*/
package opa

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"raygun/config"
	"raygun/log"
	"time"
)

/*
 *  The configuration we need to start OPA
 */
type OpaConfig struct {
	OpaPort    uint16
	OpaPath    string
	BundlePath string
	LogPath    string
}

func (oc OpaConfig) String() string {
	return fmt.Sprintf("exec: %s, bundle: %s, logs: %s", oc.OpaPath, oc.BundlePath, oc.LogPath)
}

/*
 * Details about an OPA process that is about to start, or has started
 */
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
		log.Error("Consider setting the environment variable RAYGUN_OPA_EXEC")
		return err
	}

	root_directory_logpath := fmt.Sprintf("%c%s", os.PathSeparator, config.DEFAULT_LOG_FILE)

	if opa.Config.LogPath == root_directory_logpath {
		log.Error("The environment variable TMP is not defined. OPA logs will be written to your root directory")
		log.Error("(which is almost certainly not what you want)")
		log.Error("Consider setting the environment variable TMP (to /tmp, or some equivalent on Windows)")
		log.Error("Or explicitly specifying the log path via command line arguments")
		return errors.New("invalid opa log directory: " + opa.Config.LogPath)
	}

	log.Debug("OpaRunner.Start() - commandToRun: %s - absolute_path: %s", commandToRun, absolute_path)

	args := []string{commandToRun, "run", "--server", "-b", opa.Config.BundlePath}

	log.Debug("OpaRunner.Start() - arg string: %v", args)

	opaLog, err := os.Create(opa.Config.LogPath)

	if err != nil {
		log.Error("Unable to create file: %s : %s", opa.Config.LogPath, err.Error())
		return err
	}

	process_attributes := new(os.ProcAttr)
	process_attributes.Files = []*os.File{os.Stdin, os.Stdout, opaLog}

	process, err := os.StartProcess(absolute_path, args, process_attributes)

	if err != nil {
		log.Error("Unable to start OPA: %s", err.Error())
		return err
	}

	log.Debug("Started OPA via executable: %s . Process id: %d", commandToRun, process.Pid)

	opa.Process = process

	// I don't know that we have to wait a full second for this.  TBD
	log.Debug("Waiting for 1 second for OPA to start up")
	time.Sleep(1 * time.Second)

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
