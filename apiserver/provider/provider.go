// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// Package provider defines an API end point for functions dealing
// with provider information requests.
package provider

import (
	"github.com/juju/errors"
	names "gopkg.in/juju/names.v2"

	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/facade"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/constraints"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/environs/instances"
	"github.com/juju/juju/state"
)

func init() {
	common.RegisterStandardFacade("Provider", 1, newFacade)
}

func newFacade(st *state.State, _ facade.Resources, auth facade.Authorizer) (*API, error) {
	return NewAPI(NewStateBackend(st), auth, environs.New)
}

// API publishes an endpoint to obtain provider information.
type API struct {
	backend    Backend
	modelTag   names.ModelTag
	authorizer facade.Authorizer
	environNew environs.NewEnvironFunc
}

// NewAPI returns an API pointer.
func NewAPI(backend Backend, authorizer facade.Authorizer, envNew environs.NewEnvironFunc) (*API, error) {
	if !authorizer.AuthModelManager() {
		return nil, common.ErrPerm
	}
	return &API{
		backend:    backend,
		modelTag:   backend.ModelTag(),
		authorizer: authorizer,
		environNew: envNew,
	}, nil
}

type modelEnvironConfigGetter struct {
	Backend
	model *state.Model
}

func (st *modelEnvironConfigGetter) ModelConfig() (*config.Config, error) {
	return st.model.Config()
}

func toParamsInstanceTypeResult(itypes []instances.InstanceType) []params.InstanceType {
	result := make([]params.InstanceType, len(itypes))
	for i, t := range itypes {
		virtType := ""
		if t.VirtType != nil {
			virtType = *t.VirtType
		}
		result[i] = params.InstanceType{
			Name:         t.Name,
			Arches:       t.Arches,
			CPUCores:     int(t.CpuCores),
			Memory:       int(t.Mem),
			RootDiskSize: int(t.RootDisk),
			VirtType:     virtType,
			Deprecated:   t.Deprecated,
			Cost:         int(t.Cost),
		}
	}
	return result
}

// InstanceTypes returns a list of the available instance types in the provider according
// to the passed constraints.
func (e *API) InstanceTypes(cons params.InstanceTypesConstraints) (params.InstanceTypesResults, error) {
	m, err := e.backend.GetModel(e.modelTag)
	if err != nil {
		return params.InstanceTypesResults{}, errors.Trace(err)
	}
	st := &modelEnvironConfigGetter{Backend: e.backend, model: m}
	res := make([]params.InstanceTypesResult, len(cons.Constraints))

	env, err := environs.GetEnviron(st, e.environNew)
	if err != nil {
		return params.InstanceTypesResults{}, errors.Trace(err)
	}

	for i, c := range cons.Constraints {
		consValue := constraints.Value{}
		if c.Value != nil {
			consValue = *c.Value
		}
		instanceTypes, err := env.InstanceTypes(consValue)
		allTypes := instanceTypes.InstanceTypes
		costUnit := instanceTypes.CostUnit
		costCurrency := instanceTypes.CostCurrency
		if err != nil {
			res[i] = params.InstanceTypesResult{
				Error: common.ServerError(err),
			}
			continue
		}

		res[i] = params.InstanceTypesResult{
			InstanceTypes: toParamsInstanceTypeResult(allTypes),
			CostUnit:      costUnit,
			CostCurrency:  costCurrency,
		}
	}
	return params.InstanceTypesResults{
		Results: res,
	}, nil
}
