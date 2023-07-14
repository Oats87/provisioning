#!/bin/sh

CAPI_VERSION=v1.4.4

# This script is how charts/capi/templates is generated.
# It takes the statically rendered core components and generates the corresponding chart that installs the capi-controller-manager

if ! pwd | grep -q 'charts/capi/templates'; then
  echo "Must run script from the charts/capi/templates folder"
  # Generally run in templates folder like: ../../../scripts/create-capi-chart.sh
  exit 1
fi

rm -f *capi*.yaml
rm -rf for-rancher
curl -L https://github.com/kubernetes-sigs/cluster-api/releases/download/${CAPI_VERSION}/core-components.yaml | yq -N -s '. | .kind + "-" + .metadata.name | downcase | sub("\.", "-") + ".yaml"'
mkdir for-rancher

mv customresourcedefinition-*.yaml for-rancher
mv mutatingwebhookconfiguration-capi-mutating-webhook-configuration.yaml for-rancher
mv validatingwebhookconfiguration-capi-validating-webhook-configuration.yaml for-rancher

yq -i 'del(.metadata.annotations)' for-rancher/mutatingwebhookconfiguration-capi-mutating-webhook-configuration.yaml
yq -i 'del(.metadata.labels)' for-rancher/mutatingwebhookconfiguration-capi-mutating-webhook-configuration.yaml
yq -i 'del(.metadata.name)' for-rancher/mutatingwebhookconfiguration-capi-mutating-webhook-configuration.yaml
yq -i '.metadata.creationTimestamp = null' for-rancher/mutatingwebhookconfiguration-capi-mutating-webhook-configuration.yaml
yq -i '.metadata.name = "mutating-webhook-configuration"' for-rancher/mutatingwebhookconfiguration-capi-mutating-webhook-configuration.yaml
yq -i 'del(.metadata.annotations)' for-rancher/validatingwebhookconfiguration-capi-validating-webhook-configuration.yaml
yq -i 'del(.metadata.labels)' for-rancher/validatingwebhookconfiguration-capi-validating-webhook-configuration.yaml
yq -i 'del(.metadata.name)' for-rancher/validatingwebhookconfiguration-capi-validating-webhook-configuration.yaml
yq -i '.metadata.creationTimestamp = null' for-rancher/validatingwebhookconfiguration-capi-validating-webhook-configuration.yaml
yq -i '.metadata.name = "validating-webhook-configuration"' for-rancher/validatingwebhookconfiguration-capi-validating-webhook-configuration.yaml

touch for-rancher/capi-crds.yaml
touch for-rancher/capi-webhooks.yaml
echo "---" >> for-rancher/capi-webhooks.yaml
cat for-rancher/mutatingwebhookconfiguration-capi-mutating-webhook-configuration.yaml >> for-rancher/capi-webhooks.yaml
echo "---" >> for-rancher/capi-webhooks.yaml
cat for-rancher/validatingwebhookconfiguration-capi-validating-webhook-configuration.yaml >> for-rancher/capi-webhooks.yaml
rm -f for-rancher/mutatingwebhookconfiguration-capi-mutating-webhook-configuration.yaml for-rancher/validatingwebhookconfiguration-capi-validating-webhook-configuration.yaml

for i in $(ls for-rancher/customresourcedefinition-*.yaml); do
  echo "---" >> for-rancher/capi-crds.yaml
  cat $i | yq 'del(.metadata.annotations["cert-manager.io/inject-ca-from"])' | yq 'del(.spec.conversion.webhook.clientConfig.caBundle)' >> for-rancher/capi-crds.yaml
done
sed -i '' 's/capi-system/cattle-system/g' for-rancher/*.yaml

rm -f for-rancher/customresourcedefinition-*.yaml

rm -f issuer-capi-selfsigned-issuer.yaml
rm -f certificate-capi-serving-cert.yaml
rm -f namespace-capi-system.yaml
yq -i '.metadata.annotations += {"need-a-cert.cattle.io/secret-name": "capi-webhook-service-cert"}' service-capi-webhook-service.yaml
# Remove the tolerations so we can add and customize them
yq -i 'del(.spec.template.spec.tolerations)' deployment-capi-controller-manager.yaml
yq -i 'del(.spec.template.spec.containers[0].image)' deployment-capi-controller-manager.yaml
# Remove the arguments so we can add and customize them
yq -i 'del(.spec.template.spec.containers[0].args)' deployment-capi-controller-manager.yaml
yq -i '.spec.template.spec.containers[0].image = "REPLACE-WITH-TEMPLATED-IMAGE"' deployment-capi-controller-manager.yaml
yq -i '.spec.template.spec.containers[0].env += "REPLACE-WITH-EXTRA-ENV"' deployment-capi-controller-manager.yaml
yq -i '.spec.template.spec.containers[0].args = "REPLACE-WITH-ARGS"' deployment-capi-controller-manager.yaml

cat << EOF >> deployment-capi-controller-manager.yaml
      nodeSelector: {{ include "linux-node-selector" . | nindent 8 }}
      {{- if .Values.nodeSelector }}
{{ toYaml .Values.nodeSelector | indent 8 }}
      {{- end }}
      tolerations: {{ include "linux-node-tolerations" . | nindent 6 }}
      {{- if .Values.tolerations }}
{{ toYaml .Values.tolerations | indent 6 }}
      {{- else }}
      - effect: NoSchedule
        key: node-role.kubernetes.io/controlplane
        value: "true"
      - effect: NoSchedule
        key: "node-role.kubernetes.io/control-plane"
        operator: "Exists"
      - effect: NoSchedule
        key: "node-role.kubernetes.io/master"
        operator: "Exists"
      - effect: "NoExecute"
        key: "node-role.kubernetes.io/etcd"
        operator: "Exists"
      {{- end }}
      {{- if .Values.priorityClassName }}
      priorityClassName: "{{.Values.priorityClassName}}"
      {{- end }}
EOF

sed -i '' 's/capi-system/"{{ .Release.Namespace }}"/g' *.yaml
sed -i '' 's/REPLACE-WITH-TEMPLATED-IMAGE/'\''{{ template "system_default_registry" \. }}{{ \.Values\.image\.repository }}:{{ \.Values\.image\.tag }}'\''/g' deployment-capi-controller-manager.yaml
sed -i '' 's/            - REPLACE-WITH-EXTRA-ENV/{{- if \.Values\.extraEnv }}\n{{ toYaml \.Values.extraEnv | indent 12 }}\n{{- end }}/g' deployment-capi-controller-manager.yaml
sed -i '' 's/          args: REPLACE-WITH-ARGS/          args:\n            - --leader-elect\n{{ toYaml .Values.args | indent 12 }}/g' deployment-capi-controller-manager.yaml

sed -i '' 's/IfNotPresent/"{{ .Values.image.imagePullPolicy }}"/g' deployment-capi-controller-manager.yaml

sed -i '' -e '$a\' *.yaml
