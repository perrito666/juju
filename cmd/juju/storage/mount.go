// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import (
	"github.com/juju/cmd"
	"github.com/juju/errors"
	"gopkg.in/juju/names.v2"

	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/cmd/modelcmd"
)

// NewMountVolumeCommand returns a mountVolumeCommand that implements cmd.Command.
func NewMountVolumeCommand() cmd.Command {
	cmd := &mountVolumeCommand{}
	cmd.newAPIFunc = func() (MountAPI, error) {
		return cmd.NewStorageAPI()
	}

	return modelcmd.Wrap(cmd)
}

func NewUnmountVolumeCommand() cmd.Command {
	cmd := &unmountVolumeCommand{}
	cmd.newAPIFunc = func() (UnmountAPI, error) {
		return cmd.NewStorageAPI()
	}

	return modelcmd.Wrap(cmd)
}

type baseMountCommand struct {
	StorageCommandBase

	machineTag names.MachineTag
	volumeTag  names.VolumeTag
}

const expectedArgs = 2

func (c *baseMountCommand) Init(args []string) error {
	if len(args) != expectedArgs {
		return errors.Errorf("expected %d arguments, got %d", expectedArgs, len(args))
	}
	if !names.IsValidMachine(args[0]) {
		return errors.NotValidf("machine name %q", args[0])
	}
	c.machineTag = names.NewMachineTag(args[0])

	if !names.IsValidVolume(args[1]) {
		return errors.NotValidf("volume name %q", args[1])
	}
	c.volumeTag = names.NewVolumeTag(args[1])
	return nil
}

type mountVolumeCommand struct {
	baseMountCommand

	newAPIFunc func() (MountAPI, error)
}

func (c *mountVolumeCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "mount-volume",
		Purpose: "",
		Doc:     "",
	}
}

func (c *mountVolumeCommand) Run(ctx *cmd.Context) (err error) {
	api, err := c.newAPIFunc()
	if err != nil {
		return err
	}
	defer api.Close()

	mountParams := params.MountParam{
		MachineTag: c.machineTag.String(),
		VolumeTag:  c.volumeTag.String(),
	}
	err = api.MountVolume(params.MountParams{
		MountParams: []params.MountParam{mountParams},
	})
	return err
}

// MountAPI represents an API connection that allows mounting.
type MountAPI interface {
	MountVolume(params.MountParams) error
	Close() error
}

type unmountVolumeCommand struct {
	baseMountCommand

	newAPIFunc func() (UnmountAPI, error)
}

// Info implements cmd.Command.
func (c *unmountVolumeCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "unmount-volume",
		Purpose: "",
		Doc:     "",
	}
}

func (c *unmountVolumeCommand) Run(ctx *cmd.Context) (err error) {
	api, err := c.newAPIFunc()
	if err != nil {
		return err
	}
	defer api.Close()

	unmountParams := params.MountParam{
		MachineTag: c.machineTag.String(),
		VolumeTag:  c.volumeTag.String(),
	}
	err = api.UnmountVolume(params.MountParams{
		MountParams: []params.MountParam{unmountParams},
	})
	return err
}

// UnmountAPI represents an API connection that allows unmounting.
type UnmountAPI interface {
	UnmountVolume(params.MountParams) error
	Close() error
}
