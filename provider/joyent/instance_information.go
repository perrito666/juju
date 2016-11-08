package joyent

import (
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/environs"
)

var _ environs.InstanceInformationFetcher = (*joyentEnviron)(nil)

func (env *joyentEnviron) InstanceTypes(cons params.InstanceTypesConstraints) (params.InstanceTypesResults, error) {
	result := params.InstanceTypesResults{}
	return result, nil
}
