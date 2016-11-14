// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package params

import "github.com/juju/juju/constraints"

// InstanceTypesConstraints contains a slice of InstanceTypesConstraint.
type InstanceTypesConstraints struct {
	Constraints []InstanceTypesConstraint `json:"constraints"`
}

// InstanceTypesConstraint contains a constraint applied when filtering instance types.
type InstanceTypesConstraint struct {
	// Value, if specified, contains the constraints to filter
	// the instance types by. If Value is not specified, then
	// no filtering by constraints will take place: all instance
	// types supported by the region will be returned.
	Value *constraints.Value `json:"value,omitempty"`
}

// InstanceTypesResults contains the bulk result of prompting a cloud for its instance types.
type InstanceTypesResults struct {
	Results []InstanceTypesResult `json:"results"`
}

// InstanceTypesResult contains the result of prompting a cloud for its instance types.
type InstanceTypesResult struct {
	InstanceTypes []InstanceType `json:"instance-types,omitempty"`
	CostUnit      string         `json:"cost-unit,omitempty"`
	CostCurrency  string         `json:"cost-currency,omitempty"`
	Error         *Error         `json:"error,omitempty"`
}

// InstanceType represents an available instance type in a cloud.
type InstanceType struct {
	Name         string   `json:"name,omitempty"`
	Arches       []string `json:"arches"`
	CPUCores     int      `json:"cpu-cores"`
	Memory       int      `json:"memory"`
	RootDiskSize int      `json:"root-disk,omitempty"`
	VirtType     string   `json:"virt-type,omitempty"`
	Deprecated   bool     `json:"deprecated,omitempty"`
	Cost         int      `json:"cost,omitempty"`
}
