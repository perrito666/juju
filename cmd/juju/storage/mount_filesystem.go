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

// NewMountFilesystemCommand returns a mountFilesystemCommand that implements cmd.Command.
func NewMountFilesystemCommand() cmd.Command {
	cmd := &mountFilesystemCommand{}
	cmd.newAPIFunc = func() (MountFilesystemAPI, error) {
		return cmd.NewStorageAPI()
	}

	return modelcmd.Wrap(cmd)
}

func NewUnmountFilesystemCommand() cmd.Command {
	cmd := &unmountFilesystemCommand{}
	cmd.newAPIFunc = func() (UnmountFilesystemAPI, error) {
		return cmd.NewStorageAPI()
	}

	return modelcmd.Wrap(cmd)
}

type baseMountFilesystemCommand struct {
	StorageCommandBase

	machineTag    names.MachineTag
	filesystemTag names.FilesystemTag
}

const expectedFilesystemArgs = 2

func (c *baseMountFilesystemCommand) Init(args []string) error {
	if len(args) != expectedFilesystemArgs {
		return errors.Errorf("expected %d arguments, got %d", expectedFilesystemArgs, len(args))
	}
	if !names.IsValidMachine(args[0]) {
		return errors.NotValidf("machine name %q", args[0])
	}
	c.machineTag = names.NewMachineTag(args[0])

	if !names.IsValidFilesystem(args[1]) {
		return errors.NotValidf("filesystem name %q", args[1])
	}
	c.filesystemTag = names.NewFilesystemTag(args[1])
	return nil
}

type mountFilesystemCommand struct {
	baseMountFilesystemCommand

	newAPIFunc func() (MountFilesystemAPI, error)
}

func (c *mountFilesystemCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "mount-filesystem",
		Purpose: "",
		Doc:     "",
	}
}

func (c *mountFilesystemCommand) Run(ctx *cmd.Context) (err error) {
	api, err := c.newAPIFunc()
	if err != nil {
		return err
	}
	defer api.Close()

	mountParams := params.MountFilesystemParam{
		MachineTag:    c.machineTag.String(),
		FilesystemTag: c.filesystemTag.String(),
	}
	err = api.MountFilesystem(params.MountFilesystemParams{
		MountParams: []params.MountFilesystemParam{mountParams},
	})
	return err
}

// MountFilesystemAPI represents an API connection that allows mounting.
type MountFilesystemAPI interface {
	MountFilesystem(params.MountFilesystemParams) error
	Close() error
}

type unmountFilesystemCommand struct {
	baseMountFilesystemCommand

	newAPIFunc func() (UnmountFilesystemAPI, error)
}

// Info implements cmd.Command.
func (c *unmountFilesystemCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "unmount-filesystem",
		Purpose: "",
		Doc:     "",
	}
}

func (c *unmountFilesystemCommand) Run(ctx *cmd.Context) (err error) {
	api, err := c.newAPIFunc()
	if err != nil {
		return err
	}
	defer api.Close()

	unmountParams := params.MountFilesystemParam{
		MachineTag:    c.machineTag.String(),
		FilesystemTag: c.filesystemTag.String(),
	}
	err = api.UnmountFilesystem(params.MountFilesystemParams{
		MountParams: []params.MountFilesystemParam{unmountParams},
	})
	return err
}

// UnmountFilesystemAPI represents an API connection that allows unmounting.
type UnmountFilesystemAPI interface {
	UnmountFilesystem(params.MountFilesystemParams) error
	Close() error
}
