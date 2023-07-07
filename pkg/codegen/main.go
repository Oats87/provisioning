package main

import (
	provisioningv1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"os"

	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
	v1 "k8s.io/api/core/v1"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
)

func main() {
	os.Unsetenv("GOPATH")

	controllergen.Run(args.Options{
		OutputPackage: "github.com/rancher/provisioning/pkg/generated",
		Boilerplate:   "pkg/codegen/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"provisioning.cattle.io": {
				Types: []interface{}{
					&provisioningv1.Cluster{},
				},
				GenerateTypes:   true,
				GenerateClients: true,
			},
			"rke.cattle.io": {
				Types: []interface{}{
					&rkev1.RKEControlPlane{},
					&rkev1.RKEBootstrap{},
					&rkev1.CustomMachine{},
					&rkev1.ETCDSnapshot{},
					&rkev1.RKECluster{},
					&rkev1.RKEBootstrapTemplate{},
				},
				GenerateTypes:   true,
				GenerateClients: true,
			},
			"cluster.x-k8s.io": {
				Types: []interface{}{
					capi.Machine{},
					capi.MachineSet{},
					capi.MachineDeployment{},
					capi.Cluster{},
				},
			},
			"": {
				Types: []interface{}{
					v1.Endpoints{},
					v1.PersistentVolumeClaim{},
					v1.Pod{},
					v1.Service{},
					v1.Secret{},
					v1.ConfigMap{},
					v1.ServiceAccount{},
					v1.ReplicationController{},
					v1.ResourceQuota{},
					v1.LimitRange{},
					v1.Node{},
					v1.ComponentStatus{},
					v1.Namespace{},
					v1.Event{},
				},
			},
		},
	})
}
