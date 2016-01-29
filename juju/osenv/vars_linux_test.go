// Copyright 2014 Canonical Ltd.
// Copyright 2014 Cloudbase Solutions SRL
// Licensed under the AGPLv3, see LICENCE file for details.

package osenv_test

import (
	"path/filepath"

	gc "gopkg.in/check.v1"

	"github.com/juju/juju/juju/osenv"
)

func (s *varsSuite) TestJujuData(c *gc.C) {
	path := `/foo/bar/baz/`
	// cleanup xdg config home because it has priority and it might
	// be set on the testing env.
	s.PatchEnvironment(osenv.XDGDataHome, "")
	s.PatchEnvironment("HOME", path)
	c.Assert(osenv.JujuDataLinux(), gc.Equals, filepath.Join(path, ".local", "share", "juju"))
}

func (s *varsSuite) TestJujuDataXDG(c *gc.C) {
	testJujuXDGHome := "/a/bogus/home"
	s.PatchEnvironment(osenv.XDGDataHome, testJujuXDGHome)
	homeLinux := osenv.JujuDataLinux()
	c.Assert(homeLinux, gc.Equals, filepath.Join(testJujuXDGHome, "juju"))
}

func (s *varsSuite) TestJujuDataNoXDGDefaultsConfig(c *gc.C) {
	s.PatchEnvironment(osenv.XDGDataHome, "")
	s.PatchEnvironment("HOME", "/a/bogus/user/home")
	homeLinux := osenv.JujuDataLinux()
	c.Assert(homeLinux, gc.Equals, "/a/bogus/user/home/.local/share/juju")
}
