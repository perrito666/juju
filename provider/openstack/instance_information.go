package openstack

import (
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/environs"
)

var _ environs.InstanceInformationFetcher = (*Environ)(nil)

func (e *Environ) InstanceTypes(cons params.InstanceTypesConstraints) (params.InstanceTypesResults, error) {
	result := params.InstanceTypesResults{}
	return result, nil
}
