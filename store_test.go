// Copyright 2018 Canonical Ltd.
// Licensed under the LGPL, see LICENCE file for details.

package aclstore_test

import (
	"context"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/juju/aclstore"
	"github.com/juju/simplekv/memsimplekv"
	"gopkg.in/errgo.v1"
)

func TestNewACL(t *testing.T) {
	ctx := context.Background()
	c := qt.New(t)
	store := aclstore.NewACLStore(memsimplekv.NewStore())
	err := store.NewACL(ctx, "foo", "x", "y")
	c.Assert(err, qt.Equals, nil)
	acl, err := store.Get(ctx, "foo")
	c.Assert(err, qt.Equals, nil)
	c.Assert(acl, qt.DeepEquals, []string{"x", "y"})
}

func TestNewACLOnExistingACL(t *testing.T) {
	ctx := context.Background()
	c := qt.New(t)
	store := aclstore.NewACLStore(memsimplekv.NewStore())
	err := store.NewACL(ctx, "foo", "x", "y")
	c.Assert(err, qt.Equals, nil)

	err = store.NewACL(ctx, "foo", "z", "w")
	c.Assert(err, qt.Equals, nil)

	acl, err := store.Get(ctx, "foo")
	c.Assert(err, qt.Equals, nil)
	c.Assert(acl, qt.DeepEquals, []string{"x", "y"})
}

func TestAddToNonExistentACL(t *testing.T) {
	ctx := context.Background()
	c := qt.New(t)
	store := aclstore.NewACLStore(memsimplekv.NewStore())
	err := store.Add(ctx, "foo", "x", "y")
	c.Assert(err, qt.ErrorMatches, `ACL not found`)
	c.Assert(errgo.Cause(err), qt.Equals, aclstore.ErrACLNotFound)
}

func TestAdd(t *testing.T) {
	ctx := context.Background()
	c := qt.New(t)
	store := aclstore.NewACLStore(memsimplekv.NewStore())

	err := store.NewACL(ctx, "foo", "e", "c")
	c.Assert(err, qt.Equals, nil)

	err = store.Add(ctx, "foo", "a", "d", "f", "e", "a")
	c.Assert(err, qt.Equals, nil)

	acl, err := store.Get(ctx, "foo")
	c.Assert(err, qt.Equals, nil)
	c.Assert(acl, qt.DeepEquals, []string{"a", "c", "d", "e", "f"})
}

func TestRemoveNonExistentACL(t *testing.T) {
	ctx := context.Background()
	c := qt.New(t)
	store := aclstore.NewACLStore(memsimplekv.NewStore())
	err := store.Remove(ctx, "foo", "x", "y")
	c.Assert(err, qt.ErrorMatches, `ACL not found`)
	c.Assert(errgo.Cause(err), qt.Equals, aclstore.ErrACLNotFound)
}

func TestRemove(t *testing.T) {
	ctx := context.Background()
	c := qt.New(t)
	store := aclstore.NewACLStore(memsimplekv.NewStore())

	err := store.NewACL(ctx, "foo", "a", "b", "c", "d")
	c.Assert(err, qt.Equals, nil)

	err = store.Remove(ctx, "foo", "b", "c", "e")
	c.Assert(err, qt.Equals, nil)

	acl, err := store.Get(ctx, "foo")
	c.Assert(err, qt.Equals, nil)
	c.Assert(acl, qt.DeepEquals, []string{"a", "d"})
}

func TestSetNonExistingACL(t *testing.T) {
	ctx := context.Background()
	c := qt.New(t)
	store := aclstore.NewACLStore(memsimplekv.NewStore())
	err := store.Set(ctx, "foo", "x", "y")
	c.Assert(err, qt.ErrorMatches, `ACL not found`)
	c.Assert(errgo.Cause(err), qt.Equals, aclstore.ErrACLNotFound)
}

func TestSet(t *testing.T) {
	ctx := context.Background()
	c := qt.New(t)
	store := aclstore.NewACLStore(memsimplekv.NewStore())

	err := store.NewACL(ctx, "foo", "a", "b", "c", "d")
	c.Assert(err, qt.Equals, nil)

	err = store.Set(ctx, "foo", "b", "c", "e", "e")
	c.Assert(err, qt.Equals, nil)

	acl, err := store.Get(ctx, "foo")
	c.Assert(err, qt.Equals, nil)
	c.Assert(acl, qt.DeepEquals, []string{"b", "c", "e"})
}

func TestGetNonExistent(t *testing.T) {
	ctx := context.Background()
	c := qt.New(t)
	store := aclstore.NewACLStore(memsimplekv.NewStore())

	acl, err := store.Get(ctx, "foo")
	c.Assert(err, qt.ErrorMatches, `ACL not found`)
	c.Assert(errgo.Cause(err), qt.Equals, aclstore.ErrACLNotFound)
	c.Assert(acl, qt.IsNil)
}
