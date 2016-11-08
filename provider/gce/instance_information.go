package gce

import (
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/environs"
)

var _ environs.InstanceInformationFetcher = (*environ)(nil)

func (env *environ) InstanceTypes(cons params.InstanceTypesConstraints) (params.InstanceTypesResults, error) {
	result := params.InstanceTypesResults{}
	return result, nil
}
