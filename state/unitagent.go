// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"github.com/juju/errors"
	"github.com/juju/names"
	"gopkg.in/mgo.v2/txn"
)

// UnitAgent represents the state of a service's unit agent.
type UnitAgent struct {
	st   *State
	tag  names.Tag
	name string
}

func newUnitAgent(st *State, tag names.Tag, name string) *UnitAgent {
	unitAgent := &UnitAgent{
		st:   st,
		tag:  tag,
		name: name,
	}

	return unitAgent
}

// String returns the unit agent as string.
func (u *UnitAgent) String() string {
	return u.name
}

// Status returns the status of the unit.
func (u *UnitAgent) Status() (status Status, info string, data map[string]interface{}, err error) {
	doc, err := getStatus(u.st, u.globalKey())
	if err != nil {
		return "", "", nil, errors.Trace(err)
	}
	status = doc.Status
	info = doc.StatusInfo
	data = doc.StatusData
	return
}

// SetStatus sets the status of the unit agent. The optional values
// allow to pass additional helpful status data.
func (u *UnitAgent) SetStatus(status Status, info string, data map[string]interface{}) error {
	doc, workloadDoc, err := newUnitAgentStatusDoc(status, info, data)
	if err != nil {
		return errors.Trace(err)
	}
	var ops []txn.Op
	if doc != nil {
		ops = []txn.Op{
			updateStatusOp(u.st, u.globalKey(), doc.statusDoc),
		}
	} else {
		ops = []txn.Op{
			updateStatusOp(u.st, u.globalUnitKey(), workloadDoc.statusDoc),
		}
	}
	err = u.st.runTransaction(ops)
	if err != nil {
		return errors.Errorf("cannot set status of unit agent %q: %v", u, onAbort(err, ErrDead))
	}
	return nil
}

// unitAgentGlobalKey returns the global database key for the named unit.
func unitAgentGlobalKey(name string) string {
	return "u#" + name
}

// globalKey returns the global database key for the unit.
func (u *UnitAgent) globalKey() string {
	return unitAgentGlobalKey(u.name)
}

// globalUnitKey returns the global database key for the units workload.
func (u *UnitAgent) globalUnitKey() string {
	return unitGlobalKey(u.name)
}

// Tag returns a names.Tag identifying this agent's unit.
func (u *UnitAgent) Tag() names.Tag {
	return u.tag
}
