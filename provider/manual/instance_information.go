package manual

import (
	"github.com/juju/errors"
	"github.com/juju/juju/constraints"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/instances"
)

var _ environs.InstanceTypesFetcher = (*manualEnviron)(nil)

// InstanceTypes implements InstanceTypesFetcher
func (e *manualEnviron) InstanceTypes(c constraints.Value) (instances.InstanceTypesWithCostMetadata, error) {
	result := instances.InstanceTypesWithCostMetadata{}
	return result, errors.NotSupportedf("InstanceTypes")
}
