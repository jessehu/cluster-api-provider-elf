/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"fmt"
)

type Tower struct {
	// Server is address of the tower server.
	Server string `json:"server,omitempty"`

	// Username is the name used to log into the tower server.
	Username string `json:"username,omitempty"`

	// Password is the password used to access the tower server.
	Password string `json:"password,omitempty"`
}

// ElfMachineTemplateResource describes the data needed to create a ElfMachine from a template.
type ElfMachineTemplateResource struct {
	// Spec is the specification of the desired behavior of the machine.
	Spec ElfMachineSpec `json:"spec"`
}

// APIEndpoint represents a reachable Kubernetes API endpoint.
type APIEndpoint struct {
	// The hostname on which the API server is serving.
	Host string `json:"host"`

	// The port on which the API server is serving.
	Port int32 `json:"port"`
}

// IsZero returns true if either the host or the port are zero values.
func (v APIEndpoint) IsZero() bool {
	return v.Host == "" || v.Port == 0
}

// String returns a formatted version HOST:PORT of this APIEndpoint.
func (v APIEndpoint) String() string {
	return fmt.Sprintf("%s:%d", v.Host, v.Port)
}

// NetworkStatus provides information about one of a VM's networks.
type NetworkStatus struct {
	// Connected is a flag that indicates whether this network is currently
	// connected to the VM.
	Connected bool `json:"connected,omitempty"`

	// IPAddrs is one or more IP addresses reported by vm-tools.
	// +optional
	IPAddrs []string `json:"ipAddrs,omitempty"`

	// MACAddr is the MAC address of the network device.
	MACAddr string `json:"macAddr"`

	// NetworkName is the name of the network.
	// +optional
	NetworkName string `json:"networkName,omitempty"`

	NetworkIndex int `json:"networkIndex"`
}

// NetworkSpec defines the virtual machine's network configuration.
type NetworkSpec struct {
	// Devices is the list of network devices used by the virtual machine.
	Devices []NetworkDeviceSpec `json:"devices"`

	// PreferredAPIServeCIDR is the preferred CIDR for the Kubernetes API
	// server endpoint on this machine
	PreferredAPIServerCIDR string `json:"preferredAPIServerCidr,omitempty"`

	// Vlan is the virtual LAN used by the virtual machine.
	Vlan string `json:"vlan,omitempty"`
}

// NetworkDeviceSpec defines the network configuration for a virtual machine's
// network device.
type NetworkDeviceSpec struct {
	NetworkIndex int `json:"networkIndex"`

	NetworkType string `json:"networkType"`
	// IPAddrs is a list of one or more IPv4 and/or IPv6 addresses to assign
	// to this device.
	// Required when DHCP4 and DHCP6 are both false.
	IPAddrs []string `json:"ipAddrs,omitempty"`

	Netmask string `json:"netmask,omitempty"`

	// Gateway4 is the IPv4 gateway used by this device.
	// Required when DHCP4 is false.
	Gateway string `json:"gateway,omitempty"`
}

//+kubebuilder:object:generate=false

// PatchStringValue is for patching resources.
type PatchStringValue struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}
