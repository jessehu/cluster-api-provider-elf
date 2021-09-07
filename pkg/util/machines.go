package util

import (
	"context"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "github.com/smartxworks/cluster-api-provider-elf/api/v1alpha4"
)

const (
	ProviderIDPrefix = "elf://"

	ProviderIDPattern = `(?i)^` + ProviderIDPrefix + `([a-f\d]{8}-[a-f\d]{4}-[a-f\d]{4}-[a-f\d]{4}-[a-f\d]{12})$`

	UUIDPattern = `(?i)^[a-f\d]{8}-[a-f\d]{4}-[a-f\d]{4}-[a-f\d]{4}-[a-f\d]{12}$`
)

// ErrNoMachineIPAddr indicates that no valid IP addresses were found in a machine context
var ErrNoMachineIPAddr = errors.New("no IP addresses found for machine")

// GetElfMachinesInCluster gets a cluster's ElfMachine resources.
func GetElfMachinesInCluster(
	ctx context.Context,
	controllerClient client.Client,
	namespace, clusterName string) ([]*infrav1.ElfMachine, error) {

	labels := map[string]string{clusterv1.ClusterLabelName: clusterName}
	var machineList infrav1.ElfMachineList

	if err := controllerClient.List(
		ctx, &machineList,
		client.InNamespace(namespace),
		client.MatchingLabels(labels)); err != nil {
		return nil, err
	}

	machines := make([]*infrav1.ElfMachine, len(machineList.Items))
	for i := range machineList.Items {
		machines[i] = &machineList.Items[i]
	}

	return machines, nil
}

// IsControlPlaneMachine returns true if the provided resource is
// a member of the control plane.
func IsControlPlaneMachine(machine metav1.Object) bool {
	_, ok := machine.GetLabels()[clusterv1.MachineControlPlaneLabelName]
	return ok
}

func ConvertProviderIDToUUID(providerID *string) string {
	if providerID == nil || *providerID == "" {
		return ""
	}

	pattern := regexp.MustCompile(ProviderIDPattern)
	matches := pattern.FindStringSubmatch(*providerID)
	if len(matches) < 2 {
		return ""
	}

	return matches[1]
}

func ConvertUUIDToProviderID(uuid string) string {
	if !IsUUID(uuid) {
		return ""
	}

	return ProviderIDPrefix + uuid
}

func IsUUID(uuid string) bool {
	if uuid == "" {
		return false
	}

	pattern := regexp.MustCompile(UUIDPattern)

	return pattern.MatchString(uuid)
}

func GetNetworkStatus(ipsStr string) []infrav1.NetworkStatus {
	var network []infrav1.NetworkStatus

	if ipsStr == "" {
		return network
	}

	ips := strings.Split(ipsStr, ",")
	for index, ip := range ips {
		if ip == "127.0.0.1" || strings.HasPrefix(ip, "169.254.") || strings.HasPrefix(ip, "172.17.0") {
			continue
		}

		network = append(network, infrav1.NetworkStatus{
			NetworkIndex: index,
			IPAddrs:      []string{ip},
		})
	}

	return network
}
