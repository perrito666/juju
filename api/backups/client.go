// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package backups

import (
	"github.com/juju/juju/api/base"
	"github.com/juju/juju/apiserver/params"
)

// Client wraps the backups API for the client.
type Client struct {
	base.ClientFacade
	facade base.FacadeCaller
}

// NewClient returns a new backups API client.
func NewClient(st base.APICallCloser) *Client {
	frontend, backend := base.NewClientFacade(st, "Backups")
	return &Client{ClientFacade: frontend, facade: backend}
}

// Restore is responsable for finishing a restore after a placeholder
// machine has been bootstraped, it receives the name of a backup
// file on server and will return error on failure.
func (c *Client) Restore(backupFileName, backupId string) error {
	params := params.Restore{FileName: backupFileName,
		BackupId: backupId,
		Machine:  "0"}
	err := c.facade.FacadeCall("Restore", params, nil)
	return err
}

func (c *Client) PrepareRestore() error {
	return c.facade.FacadeCall("PrepareRestore", nil, nil)
}

func (c *Client) FinishRestore() error {
	return c.facade.FacadeCall("FinishRestore", nil, nil)
}

// PublicAddress returns the public address of the specified
// machine or unit. For a machine, target is an id not a tag.
func (c *Client) PublicAddress(target string) (string, error) {
	var results params.PublicAddressResults
	p := params.PublicAddress{Target: target}
	err := c.facade.FacadeCall("PublicAddress", p, &results)
	return results.PublicAddress, err
}
