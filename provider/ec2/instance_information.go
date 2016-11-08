package ec2

import (
	"github.com/juju/errors"

	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/environs"
)

var _ environs.InstanceInformationFetcher = (*environ)(nil)

func (e *environ) InstanceTypes(cons params.InstanceTypesConstraints) (params.InstanceTypesResults, error) {
	iTypes, err := e.supportedInstanceTypes()
	if err != nil {
		return params.InstanceTypesResults{}, errors.Trace(err)
	}
	result := params.InstanceTypesResults{}
	result.Results = make([]params.InstanceTypesResult, len(iTypes))
	for i, t := range iTypes {
		res := params.InstanceTypesResult{
			InstanceTypes: t.ToParamsInstanceType(),
			CostUnit:      "TODO",
			CostCurrency:  "USD",
		}
		result.Results[i] = res
	}
	return result, nil
}
