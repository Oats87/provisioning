package capr

import (
	"context"

	"github.com/rancher/provisioning/capr"
	"github.com/rancher/provisioning/capr/planner"
	"github.com/rancher/provisioning/controllers/capr/bootstrap"
	"github.com/rancher/provisioning/controllers/capr/dynamicschema"
	"github.com/rancher/provisioning/controllers/capr/machinedrain"
	"github.com/rancher/provisioning/controllers/capr/machinenodelookup"
	"github.com/rancher/provisioning/controllers/capr/machineprovision"
	"github.com/rancher/provisioning/controllers/capr/managesystemagent"
	plannercontroller "github.com/rancher/provisioning/controllers/capr/planner"
	"github.com/rancher/provisioning/controllers/capr/plansecret"
	"github.com/rancher/provisioning/controllers/capr/rkecluster"
	"github.com/rancher/provisioning/controllers/capr/rkecontrolplane"
	"github.com/rancher/provisioning/controllers/capr/unmanaged"
	"github.com/rancher/provisioning/provisioningv2/image"
	"github.com/rancher/provisioning/provisioningv2/kubeconfig"
	"github.com/rancher/provisioning/provisioningv2/systeminfo"
	"github.com/rancher/provisioning/wrangler"
	"github.com/rancher/rancher/pkg/settings"
)

func Register(ctx context.Context, clients *wrangler.Context, kubeconfigManager *kubeconfig.Manager) {
	rkePlanner := planner.New(ctx, clients, planner.InfoFunctions{
		ImageResolver:           image.ResolveWithControlPlane,
		ReleaseData:             capr.GetKDMReleaseData,
		SystemAgentImage:        settings.SystemAgentInstallerImage.Get,
		SystemPodLabelSelectors: systeminfo.NewRetriever(clients).GetSystemPodLabelSelectors,
	})
	//if features.MCM.Enabled() {
	dynamicschema.Register(ctx, clients)
	machineprovision.Register(ctx, clients, kubeconfigManager)
	//}
	rkecluster.Register(ctx, clients)
	bootstrap.Register(ctx, clients)
	machinenodelookup.Register(ctx, clients, kubeconfigManager)
	plannercontroller.Register(ctx, clients, rkePlanner)
	plansecret.Register(ctx, clients)
	unmanaged.Register(ctx, clients, kubeconfigManager)
	rkecontrolplane.Register(ctx, clients)
	managesystemagent.Register(ctx, clients)
	machinedrain.Register(ctx, clients)
}
