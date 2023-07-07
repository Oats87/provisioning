package main

import (
	"context"
	"flag"
	"github.com/ehazlett/simplelog"
	"github.com/rancher/provisioning/controllers/capi"
	"github.com/rancher/provisioning/controllers/capr"
	"github.com/rancher/provisioning/controllers/provisioningv2"
	"github.com/rancher/provisioning/controllers/provisioningv2/cluster"
	"github.com/rancher/provisioning/crds"
	"github.com/rancher/provisioning/provisioningv2/kubeconfig"
	"github.com/rancher/provisioning/wrangler"
	"github.com/rancher/rancher/pkg/controllers/dashboardapi/settings"
	wKubeconfig "github.com/rancher/wrangler/pkg/kubeconfig"
	"github.com/rancher/wrangler/pkg/signals"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	masterURL      string
	kubeconfigFile string
)

func init() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.StringVar(&kubeconfigFile, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.Parse()
}

func main() {
	logrus.SetFormatter(&simplelog.StandardFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)
	// set up signals so we handle the first shutdown signal gracefully
	ctx := signals.SetupSignalContext()

	// This will load the kubeconfig file in a style the same as kubectl
	cfg := wKubeconfig.GetNonInteractiveClientConfig(kubeconfigFile)

	clientConfig, err := cfg.ClientConfig()
	if err != nil {
		logrus.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	wranglerContext, err := wrangler.NewContext(ctx, cfg, clientConfig)
	if err != nil {
		logrus.Fatalf("Error building wrangler context: %s", err.Error())
	}

	cluster.RegisterIndexers(wranglerContext)
	if err = crds.Create(ctx, clientConfig); err != nil {
		logrus.Fatalf("Error creating CRDs: %s", err.Error())
	}

	// Register settings so that the provider is set and we can retrieve the internal server URL + CA for the kubeconfig manager below.
	if err = settings.Register(wranglerContext.Mgmt.Setting()); err != nil {
		logrus.Fatalf("Error registering settings: %s", err.Error())
	}

	kubeconfigManager := kubeconfig.New(wranglerContext)
	provisioningv2.Register(ctx, wranglerContext, kubeconfigManager)
	capr.Register(ctx, wranglerContext, kubeconfigManager)

	capiStart, err := capi.Register(ctx, wranglerContext)
	if err != nil {
		logrus.Fatalf("Error building CAPI controllers: %s", err.Error())
	}
	wranglerContext.OnLeader(func(ctx context.Context) error {
		if err := capiStart(ctx); err != nil {
			logrus.Fatal(err)
		}
		logrus.Info("Cluster API is started")
		return nil
	})

	wranglerContext.Start(ctx)

	<-ctx.Done()
}
