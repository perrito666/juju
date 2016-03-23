// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package joyent

import (
	gc "gopkg.in/check.v1"

	coretesting "github.com/juju/juju/testing"
)

type InternalSuite struct {
	coretesting.FakeJujuXDGDataHomeSuite
}

var _ = gc.Suite(&InternalSuite{})

func (s *InternalSuite) TestEnsurePrivateKey(c *gc.C) {
	m := map[string]interface{}{
		"private-key": "foo",
	}

	e := &environConfig{attrs: copymap(m)}

	err := ensurePrivateKeyOrPath(e)
	c.Assert(err, gc.IsNil)
	c.Assert(e.attrs, gc.DeepEquals, m)
}

func (s *InternalSuite) TestEnsurePrivateKeyPathSet(c *gc.C) {
	// if path is set, func should noop
	m := map[string]interface{}{
		"private-key-path": "foo",
	}

	e := &environConfig{attrs: copymap(m)}

	err := ensurePrivateKeyOrPath(e)
	c.Assert(err, gc.IsNil)
	c.Assert(e.attrs, gc.DeepEquals, m)
}

func (s *InternalSuite) TestEnsurePrivateKeySet(c *gc.C) {
	// if key is set, func should noop
	m := map[string]interface{}{"private-key": "foo"}

	e := &environConfig{attrs: copymap(m)}

	err := ensurePrivateKeyOrPath(e)
	c.Assert(err, gc.IsNil)
	c.Assert(e.attrs, gc.DeepEquals, m)
}

func (s *InternalSuite) TestEnsurePrivateKeyMissing(c *gc.C) {
	e := &environConfig{attrs: map[string]interface{}{}}

	err := ensurePrivateKeyOrPath(e)
	c.Assert(err, gc.ErrorMatches, "no ssh private key specified in joyent configuration")
}

func copymap(m map[string]interface{}) map[string]interface{} {
	m1 := make(map[string]interface{}, len(m))
	for k, v := range m {
		m1[k] = v
	}
	return m1
}
