package maas

import (
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/environs"
)

var _ environs.InstanceInformationFetcher = (*maasEnviron)(nil)

func (env *maasEnviron) InstanceTypes(cons params.InstanceTypesConstraints) (params.InstanceTypesResults, error) {
	result := params.InstanceTypesResults{}
	return result, nil
}
