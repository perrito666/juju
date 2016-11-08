package azure

import (
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/environs"
)

var _ environs.InstanceInformationFetcher = (*azureEnviron)(nil)

func (env *azureEnviron) InstanceTypes(cons params.InstanceTypesConstraints) (params.InstanceTypesResults, error) {
	result := params.InstanceTypesResults{}
	return result, nil
}
