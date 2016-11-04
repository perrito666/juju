// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package provider

import (
	names "gopkg.in/juju/names.v2"

	"github.com/juju/juju/environs"
	"github.com/juju/juju/state"
	"github.com/juju/juju/state/stateenvirons"
)

// Backend provides the methods necessary to obtain provider
// information.
type Backend interface {
	environs.EnvironConfigGetter
	ModelTag() names.ModelTag
	Close() error
	GetModel(names.ModelTag) (*state.Model, error)
}

// NewStateBackend returns a Backend that wraps the passed state.
func NewStateBackend(st *state.State) Backend {
	return stateenvirons.EnvironConfigGetter{st}
}
