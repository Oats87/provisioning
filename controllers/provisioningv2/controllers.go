package provisioningv2

import (
	"context"

	"github.com/rancher/provisioning/controllers/provisioningv2/cluster"
	"github.com/rancher/provisioning/controllers/provisioningv2/clusterindex"
	"github.com/rancher/provisioning/controllers/provisioningv2/provisioningcluster"
	"github.com/rancher/provisioning/controllers/provisioningv2/provisioninglog"
	"github.com/rancher/provisioning/controllers/provisioningv2/secret"
	"github.com/rancher/provisioning/provisioningv2/kubeconfig"
	"github.com/rancher/provisioning/wrangler"
)

func Register(ctx context.Context, clients *wrangler.Context, kubeconfigManager *kubeconfig.Manager) {
	cluster.Register(ctx, clients, kubeconfigManager)
	clusterindex.Register(ctx, clients)
	//if features.MCM.Enabled() {
	secret.Register(ctx, clients)
	//}
	provisioningcluster.Register(ctx, clients)
	provisioninglog.Register(ctx, clients)

	//if features.Fleet.Enabled() {
	//managedchart.Register(ctx, clients)
	//fleetcluster.Register(ctx, clients)
	//fleetworkspace.Register(ctx, clients)
	//}
}
