// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import (
	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/params"
	"gopkg.in/juju/names.v2"
)

func parseMountTags(machine, volume string) (names.MachineTag, names.VolumeTag, error) {
	machineTag, err := names.ParseMachineTag(machine)
	if err != nil {
		return names.MachineTag{}, names.VolumeTag{}, err
	}
	volumeTag, err := names.ParseVolumeTag(volume)
	if err != nil {
		return names.MachineTag{}, names.VolumeTag{}, err
	}
	return machineTag, volumeTag, err
}

// MountVolume will attach a volume identified by mountParams.VolumeTag to a
// machine identified by mountParams.MachineTag.
func (api *API) MountVolume(mountParams params.MountParams) error {
	for _, mountParam := range mountParams.MountParams {
		if err := mountVolume(api, mountParam); err != nil {
			return common.ServerError(err)
		}
	}
	return nil
}

func mountVolume(api *API, mountParam params.MountParam) error {
	machineTag, volumeTag, err := parseMountTags(mountParam.MachineTag, mountParam.VolumeTag)
	if err != nil {
		return err
	}
	if err := api.storage.AttachVolume(machineTag, volumeTag); err != nil {
		return err
	}
	return nil

}

// UnmountVolume will detach a volume identified by unmountParams.VolumeTag from a
// machine identified by unmountParams.MachineTag or return error if not possible.
func (api *API) UnmountVolume(unmountParams params.MountParams) error {
	for _, unmountParam := range unmountParams.MountParams {
		if err := unmountVolume(api, unmountParam); err != nil {
			return common.ServerError(err)
		}
	}
	return nil
}

func unmountVolume(api *API, unmountParam params.MountParam) error {
	machineTag, volumeTag, err := parseMountTags(unmountParam.MachineTag, unmountParam.VolumeTag)
	if err != nil {
		return err
	}
	if err := api.storage.DetachVolume(machineTag, volumeTag); err != nil {
		return err
	}
	return nil
}
