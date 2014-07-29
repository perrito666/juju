// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/juju/cmd"
	"launchpad.net/gnuflag"

	"github.com/juju/juju/cmd/envcmd"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/bootstrap"
	"github.com/juju/juju/environs/configstore"
	"github.com/juju/juju/constraints"
	"github.com/juju/juju/juju"
	"github.com/juju/juju/provider/common"
	"github.com/juju/juju/state/api/params"
)

type restoreClient interface {
	Restore(io.Reader) error
}

type RestoreCommand struct {
	envcmd.EnvCommandBase
	Constraints constraints.Value
	Filename    string
	Upload	bool
	Bootstrap bool
}

var restoreDoc = `
Restores a backup that was previously created with "juju backup".

This command creates a new state server and arranges for it to replace
the previous state server for an environment.  It does *not* restore
an existing server to a previous state, but instead creates a new server
with equivanlent state.  As part of restore, all known instances are
configured to treat the new state server as their master.

The given constraints will be used to choose the new instance.

If the provided state cannot be restored, this command will fail with
an appropriate message.  For instance, if the existing bootstrap
instance is already running then the command will fail with a message
to that effect.
`

func (c *RestoreCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "restore",
		Purpose: "restore a state server backup made with juju backup",
		Args:    "[-u] [-b] <backupfile.tar.gz>",
		Doc:     strings.TrimSpace(restoreDoc),
	}
}

func (c *RestoreCommand) SetFlags(f *gnuflag.FlagSet) {
	f.Var(constraints.ConstraintsValue{Target: &c.Constraints},
		"constraints", "set environment constraints")
}

func (c *RestoreCommand) Init(args []string) error {
	for _, arg := range args {
		switch {
		case arg == "-u":
			c.Upload = true
		case arg == "-b":
			c.Bootstrap = true
		case c.Filename == "":
			c.Filename = arg
		default:
			return fmt.Errorf("unrecognized argument: %v", arg)
	}
	if c.Filename == "" {
		return fmt.Errorf("no backup name specified")
	} 
	return nil
}

const restoreAPIIncompatibility = "server version not compatible for " +
	"restore with client version"

func (c *RestoreCommand) runRestore(ctx *cmd.Context, client restoreClient) error {
	fd, err := os.Open(c.Filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	fileName := filepath.Base(c.Filename)
	if err := client.Restore(fileName); err != nil {
		if params.IsCodeNotImplemented(err) {
			return fmt.Errorf(restoreAPIIncompatibility)
		} 
		return err
	}
	fmt.Fprintf(ctx.Stdout, "restore from %s completed\n", c.Filename)
	return nil
}

func (c *RestoreCommand) rebootstrap (ctx *cmd.Context) (environs.Environ, error) {
	cons := c.Constraints
	store, err := configstore.Default()
	if err != nil {
		return nil, err
	}
	cfg, err := c.Config(store)
	if err != nil {
		return nil, err
	}
	// Turn on safe mode so that the newly bootstrapped instance
	// will not destroy all the instances it does not know about.
	cfg, err = cfg.Apply(map[string]interface{}{
		"provisioner-safe-mode": true,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot enable provisioner-safe-mode: %v", err)
	}
	env, err := environs.New(cfg)
	if err != nil {
		return nil, err
	}
	instanceIds, err := env.StateServerInstances()
	if err != nil {
		return nil, fmt.Errorf("cannot determine state server instances: %v", err)
	}
	if len(instanceIds) == 0 {
		return nil, fmt.Errorf("no instances found; perhaps the environment was not bootstrapped")
	}
	if len(instanceIds) > 1 {
		return nil, fmt.Errorf("restore does not support HA juju configurations yet")
	}
	inst, err := env.Instances(instanceIds)
	if err == nil {
		return nil, fmt.Errorf("old bootstrap instance %q still seems to exist; will not replace", inst)
	}
	if err != environs.ErrNoInstances {
		return nil, fmt.Errorf("cannot detect whether old instance is still running: %v", err)
	}
	// Remove the storage so that we can bootstrap without the provider complaining.
	if err := env.Storage().Remove(common.StateFile); err != nil {
		return nil, fmt.Errorf("cannot remove %q from storage: %v", common.StateFile, err)
	}

	// TODO If we fail beyond here, then we won't have a state file and
	// we won't be able to re-run this script because it fails without it.
	// We could either try to recreate the file if we fail (which is itself
	// error-prone) or we could provide a --no-check flag to make
	// it go ahead anyway without the check.

	args := environs.BootstrapParams{Constraints: cons}
	if err := bootstrap.Bootstrap(ctx, env, args); err != nil {
		return nil, fmt.Errorf("cannot bootstrap new instance: %v", err)
	}
	return env, nil
}

func (c *RestoreCommand) doUpload(client *juju.Client) error {	
	return nil
}

func (c *RestoreCommand) Run(ctx *cmd.Context) error {
	if c.Bootstrap {
		_, err := c.rebootstrap(ctx)
		if err != nil {
			return err
		}
	}
	// Empty string will get a client for current default
	client, err := juju.NewAPIClientFromName("")
	if err != nil {
		return err
	}
	defer client.Close()

	if c.Upload {
		c.doUpload(client)
	}

	return c.runRestore(ctx, client)
}

