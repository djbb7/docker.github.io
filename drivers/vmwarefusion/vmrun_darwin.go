/*
 * Copyright 2014 VMware, Inc.  All rights reserved.  Licensed under the Apache v2 License.
 */

package vmwarefusion

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/docker/machine/log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	vmrunbin    = setVmwareCmd("vmrun")
	vdiskmanbin = setVmwareCmd("vmware-vdiskmanager")
)

var (
	ErrMachineExist    = errors.New("machine already exists")
	ErrMachineNotExist = errors.New("machine does not exist")
	ErrVMRUNNotFound   = errors.New("VMRUN not found")
)

// detect the vmrun and vmware-vdiskmanager cmds' path if needed
func setVmwareCmd(cmd string) string {
	if path, err := exec.LookPath(cmd); err == nil {
		return path
	}
	return filepath.Join("/Applications/VMware Fusion.app/Contents/Library/", cmd)
}

func vmrun(args ...string) (string, string, error) {
	cmd := exec.Command(vmrunbin, args...)
	if os.Getenv("DEBUG") != "" {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	log.Debugf("executing: %v %v", vmrunbin, strings.Join(args, " "))

	err := cmd.Run()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			err = ErrVMRUNNotFound
		}
	}

	return stdout.String(), stderr.String(), err
}

// Make a vmdk disk image with the given size (in MB).
func vdiskmanager(dest string, size int) error {
	cmd := exec.Command(vdiskmanbin, "-c", "-t", "0", "-s", fmt.Sprintf("%dMB", size), "-a", "lsilogic", dest)
	if os.Getenv("DEBUG") != "" {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if stdout := cmd.Run(); stdout != nil {
		if ee, ok := stdout.(*exec.Error); ok && ee == exec.ErrNotFound {
			return ErrVMRUNNotFound
		}
	}
	return nil
}
