package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("rkt-oci", "Open Container Initiative wrapper for rkt")

	// centralized error handling
	var err error
	app.After = func() {
		if err != nil {
			logrus.Println(err.Error())
			cli.Exit(1)
		}
	}

	app.Command(
		"version",
		"Print the version and exit",
		func(cmd *cli.Cmd) {
			cmd.Action = func() {
				err = version()
			}
		})
	app.Command(
		"start",
		"Create a container and launch a process inside it",
		func(cmd *cli.Cmd) {
			cmd.Action = func() {
				err = start()
			}
		})

	app.Run(os.Args)
}

func version() error {
	_, err := fmt.Println("rkt-oci TODO")
	return err
}

func convert(appName string, workingDir string) error {
	var cmd *exec.Cmd
	aciName := appName + ".aci"
	//set appName to rkt appname, set rkt aciName to image name
	cmd = exec.Command("oci2aci", "--debug", "--name", appName, appName, aciName)
	cmd.Dir = workingDir
	out, err := cmd.CombinedOutput()

	if err != nil {
		logrus.Debugf(string(out))
		return err
	}

	return nil
}

func start() error {
	logrus.Debugf("starting container")

	specDir, err := os.Getwd()
	if err != nil {
		return err
	}
	appName := filepath.Base(specDir)
	aciName := appName + ".aci"
	aciPath := filepath.Dir(specDir)

	if err = convert(appName, aciPath); err != nil {
		return err
	}

	cmd := exec.Command("rkt", "run", aciName, "--interactive", "--insecure-skip-verify", "--mds"+
		"-register=false", "--volume", "proc,kind=host,source=/bin", "--volume", "dev,kind=host,"+
		"source=/bin", "--volume", "devpts,kind=host,source=/bin", "--volume", "shm,kind=host,"+
		"source=/bin", "--volume", "mqueue,kind=host,source=/bin", "--volume", "sysfs,kind=host,"+
		"source=/bin", "--volume", "cgroup,kind=host,source=/bin", "--net=host")
	cmd.Dir = aciPath
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
