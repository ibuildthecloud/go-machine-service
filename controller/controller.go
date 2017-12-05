package controller

import (
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
)

func Register(management *config.ManagementContext) {
	machineLifecycle := &MachineLifecycle{}
	machineClient := management.Management.Machines("")

	machineClient.
		Controller().
		AddHandler(v3.NewMachineLifecycleAdapter("machine-controller", machineClient, machineLifecycle))
}

type MachineLifecycle struct {
}

func (m *MachineLifecycle) Create(obj *v3.Machine) error {
	// No need to create a deepcopy of obj, obj is already a deepcopy
	return nil
}

func (m *MachineLifecycle) Updated(obj *v3.Machine) error {
	// YOU MUST CALL DEEPCOPY
	return nil
}

func (m *MachineLifecycle) Remove(obj *v3.Machine) error {
	// No need to create a deepcopy of obj, obj is already a deepcopy
	return nil
}
