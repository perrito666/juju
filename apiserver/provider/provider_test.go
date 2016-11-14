// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package provider

import (
	"github.com/juju/errors"
	jujutesting "github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	names "gopkg.in/juju/names.v2"

	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/apiserver/testing"
	"github.com/juju/juju/constraints"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/environs/instances"
	"github.com/juju/juju/provider/dummy"
	"github.com/juju/juju/state"
)

type providerTypesSuite struct{}

var _ = gc.Suite(&providerTypesSuite{})

var over9kCPUCores uint64 = 9001

func (p *providerTypesSuite) TestInstanceTypes(c *gc.C) {
	backend := mockBackend{}
	authorizer := testing.FakeAuthorizer{Tag: names.NewUserTag("admin"),
		EnvironManager: true}
	itCons := constraints.Value{CpuCores: &over9kCPUCores}
	failureCons := constraints.Value{}
	m := mockEnviron{
		results: map[constraints.Value]instances.InstanceTypesWithCostMetadata{
			itCons: instances.InstanceTypesWithCostMetadata{
				CostUnit:     "USD/h",
				CostCurrency: "USD",
				InstanceTypes: []instances.InstanceType{
					{Name: "instancetype-1"},
					{Name: "instancetype-2"}},
			},
		},
	}
	newMockEnviron := func(args environs.OpenParams) (environs.Environ, error) {
		return &m, nil
	}
	api := API{
		backend:              &backend,
		modelTag:             backend.ModelTag(),
		authorizer:           authorizer,
		environNew:           newMockEnviron,
		newModelConfigGetter: newMockModelConfigGetter,
	}

	cons := params.InstanceTypesConstraints{
		Constraints: []params.InstanceTypesConstraint{{Value: &itCons}, {Value: &failureCons}},
	}
	r, err := api.InstanceTypes(cons)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(r.Results, gc.HasLen, 2)
	expected := []params.InstanceTypesResult{
		params.InstanceTypesResult{
			InstanceTypes: []params.InstanceType{
				params.InstanceType{
					Name: "instancetype-1",
				},
				params.InstanceType{
					Name: "instancetype-2",
				}},
			CostUnit:     "USD/h",
			CostCurrency: "USD",
		},
		params.InstanceTypesResult{
			Error: &params.Error{Message: "Instances matching constraint  not found", Code: "not found"}}}
	c.Assert(r.Results, gc.DeepEquals, expected)
}

type mockBackend struct {
	jujutesting.Stub
	Backend

	Cloud environs.CloudSpec
}

func (*mockBackend) ModelTag() names.ModelTag {
	return names.NewModelTag("beef1beef1-0000-0000-000011112222")
}

func (*mockBackend) GetModel(t names.ModelTag) (*state.Model, error) {
	return &state.Model{}, nil
}

func (fb *mockBackend) CloudSpec(names.ModelTag) (environs.CloudSpec, error) {
	fb.MethodCall(fb, "CloudSpec")
	if err := fb.NextErr(); err != nil {
		return environs.CloudSpec{}, err
	}
	return fb.Cloud, nil
}

type mockEnviron struct {
	environs.Environ
	jujutesting.Stub

	results map[constraints.Value]instances.InstanceTypesWithCostMetadata
}

func (m *mockEnviron) InstanceTypes(c constraints.Value) (instances.InstanceTypesWithCostMetadata, error) {
	it, ok := m.results[c]
	if !ok {
		return instances.InstanceTypesWithCostMetadata{}, errors.NotFoundf("Instances matching constraint %v", c)
	}
	return it, nil
}

func newMockModelConfigGetter(b Backend, m modelConfigAble) *modelEnvironConfigGetter {
	return &modelEnvironConfigGetter{Backend: b, model: &mockModel{}}
}

type mockModel struct {
}

func (*mockModel) Config() (*config.Config, error) {
	return config.New(config.UseDefaults, dummy.SampleConfig())
}
