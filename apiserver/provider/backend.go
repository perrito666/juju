package provider

import (
	"github.com/juju/juju/environs"
	"github.com/juju/juju/state"
	names "gopkg.in/juju/names.v2"
)

type Backend interface {
	environs.EnvironConfigGetter
	Close() error
	GetModel(names.ModelTag) (*state.Model, error)
}

type stateShim struct {
	*state.State
}

func NewStateBackend(st *state.State) Backend {
	return stateShim{st}
}

func (s stateShim) ControllerModel() (Model, error) {
	m, err := s.State.ControllerModel()
	if err != nil {
		return nil, err
	}
	return m, nil
}
