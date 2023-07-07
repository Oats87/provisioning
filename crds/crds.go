package crds

import (
	"context"

	"github.com/rancher/provisioning/crds/provisioningv2"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/crd"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

func List(cfg *rest.Config) (_ []crd.CRD, err error) {
	var result []crd.CRD

	//if features.ProvisioningV2.Enabled() {
	result = []crd.CRD{crd.CRD{
		SchemaObject: v3.ManagedChart{},
	}.WithStatus(),
	}
	//}

	//if features.ProvisioningV2.Enabled() {
	result = append(result, provisioningv2.List()...)
	//}
	return result, nil
}

func Webhooks() []runtime.Object {
	//if features.ProvisioningV2.Enabled() {
	return provisioningv2.Webhooks()
	//}
	//return nil
}

func Create(ctx context.Context, cfg *rest.Config) error {
	factory, err := crd.NewFactoryFromClient(cfg)
	if err != nil {
		return err
	}

	apply, err := apply.NewForConfig(cfg)
	if err != nil {
		return err
	}
	apply = apply.
		WithSetID("crd-webhooks").
		WithDynamicLookup().
		WithNoDelete()
	if err := apply.ApplyObjects(Webhooks()...); err != nil {
		return err
	}

	crds, err := List(cfg)
	if err != nil {
		return err
	}

	return factory.BatchCreateCRDs(ctx, crds...).BatchWait()
}

func newCRD(obj interface{}, customize func(crd.CRD) crd.CRD) crd.CRD {
	crd := crd.CRD{
		SchemaObject: obj,
	}
	if customize != nil {
		crd = customize(crd)
	}
	return crd
}
