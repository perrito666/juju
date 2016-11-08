package params

import "github.com/juju/juju/constraints"

type InstanceTypesConstraints struct {
	Constraints []InstanceTypesConstraint `json:"constraints"`
}

type InstanceTypesConstraint struct {
	// Value, if specified, contains the constraints to filter
	// the instance types by. If Value is not specified, then
	// no filtering by constraints will take place: all instance
	// types supported by the region will be returned.
	Value *constraints.Value `json:"value,omitempty"`
}

type InstanceTypesResults struct {
	Results []InstanceTypesResult `json:"results"`
}

type InstanceTypesResult struct {
	InstanceTypes []InstanceType `json:"instance-types,omitempty"`
	CostUnit      string         `json:"cost-unit,omitempty"`
	CostCurrency  string         `json:"cost-currency,omitempty"`
	Error         *Error         `json:"error,omitempty"`
}

type InstanceType struct {
	Name         string `json:"name,omitempty"`
	Arch         string `json:"arch"`
	CPUCores     int    `json:"cpu-cores"`
	Memory       int    `json:"memory"`
	RootDiskSize int    `json:"root-disk,omitempty"`
	VirtType     string `json:"virt-type,omitempty"`
	Deprecated   bool   `json:"deprecated,omitempty"`
	Cost         int    `json:"cost,omitempty"`
}
