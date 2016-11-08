// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// Package provider defines an API end point for functions dealing
// with provider information requests.
package provider

import (
	"github.com/juju/errors"
	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/facade"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/state"
	names "gopkg.in/juju/names.v2"
)

func init() {
	common.RegisterStandardFacade("Provider", 1, newFacade)
}

func newFacade(st *state.State, _ facade.Resources, auth facade.Authorizer) (*API, error) {
	return NewAPI(NewStateBackend(st), auth)
}

// API publishes an endpoint to obtain provider information.
type API struct {
	backend    Backend
	authorizer facade.Authorizer
}

// NewAPI returns an API pointer.
func NewAPI(backend Backend, authorizer facade.Authorizer) (*API, error) {
	return &API{backend: backend, authorizer: authorizer}, nil
}

// InstanceTypes returns a list of the available instance types in the provider according
// to the passed constraints.
func (e *API) InstanceTypes(uuid string, cons params.InstanceTypesConstraints) (params.InstanceTypesResults, error) {
	if !names.IsValidModel(uuid) {
		return params.InstanceTypesResults{}, errors.NotValidf("%q uuid", uuid)
	}
	modelTag := names.NewModelTag(uuid)
	m, err := e.backend.GetModel(modelTag)
	if err != nil {
		return params.InstanceTypesResults{}, errors.Trace(err)
	}
	// Perhaps arrange stateshim to return this Model when GetModelConfig is invoked?
	environs.GetEnviron(e.backend, environs.New)
}
