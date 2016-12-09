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
	cmd.newAPIFunc = func() (MountVolumeAPI, error) {
		return cmd.NewStorageAPI()
	}

	return modelcmd.Wrap(cmd)
}

func NewUnmountVolumeCommand() cmd.Command {
	cmd := &unmountVolumeCommand{}
	cmd.newAPIFunc = func() (UnmountVolumeAPI, error) {
		return cmd.NewStorageAPI()
	}

	return modelcmd.Wrap(cmd)
}

type baseMountVolumeCommand struct {
	StorageCommandBase

	machineTag names.MachineTag
	volumeTag  names.VolumeTag
}

const expectedVolumeArgs = 2

func (c *baseMountVolumeCommand) Init(args []string) error {
	if len(args) != expectedVolumeArgs {
		return errors.Errorf("expected %d arguments, got %d", expectedVolumeArgs, len(args))
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
	baseMountVolumeCommand

	newAPIFunc func() (MountVolumeAPI, error)
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

	mountParams := params.MountVolumeParam{
		MachineTag: c.machineTag.String(),
		VolumeTag:  c.volumeTag.String(),
	}
	err = api.MountVolume(params.MountVolumeParams{
		MountParams: []params.MountVolumeParam{mountParams},
	})
	return err
}

// MountVolumeAPI represents an API connection that allows mounting.
type MountVolumeAPI interface {
	MountVolume(params.MountVolumeParams) error
	Close() error
}

type unmountVolumeCommand struct {
	baseMountVolumeCommand

	newAPIFunc func() (UnmountVolumeAPI, error)
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

	unmountParams := params.MountVolumeParam{
		MachineTag: c.machineTag.String(),
		VolumeTag:  c.volumeTag.String(),
	}
	err = api.UnmountVolume(params.MountVolumeParams{
		MountParams: []params.MountVolumeParam{unmountParams},
	})
	return err
}

// UnmountVolumeAPI represents an API connection that allows unmounting.
type UnmountVolumeAPI interface {
	UnmountVolume(params.MountVolumeParams) error
	Close() error
}
