#!/usr/bin/env bash
#
set -x

function start-docker-desktop {
  open -a Docker
  sleep 20
}

function install-cluster-no-kubeproxy {
echo "Installing 2x worker - 1x master Kind with Ingress ports exposed..."
kind delete cluster
kind create cluster --config - <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  disableDefaultCNI: true
  kubeProxyMode: none
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    listenAddress: 127.0.0.1
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    listenAddress: 127.0.0.1
    protocol: TCP
- role: worker
- role: worker
EOF
}

function install-cilium {
echo "Installing cilium with hubble and default kind-control-plane node..."
echo "'docker network kind inspect' for more info..."
helm upgrade --install --namespace kube-system --repo https://helm.cilium.io cilium cilium --values - <<EOF
kubeProxyReplacement: strict
k8sServiceHost: kind-control-plane
k8sServicePort: 6443
hostServices:
  enabled: false
externalIPs:
  enabled: true
nodePort:
  enabled: true
hostPort:
  enabled: true
image:
  pullPolicy: IfNotPresent
ipam:
  mode: kubernetes
hubble:
  enabled: true
  relay:
    enabled: true
  ui:
    enabled: true
    ingress:
      enabled: true
      annotations:
        kubernetes.io/ingress.class: nginx
      hosts:
        - hubble-ui.127.0.0.1.nip.io
EOF
}

function install-ingress-nginx {
echo "Installing ingress NGINX..."
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.1.1/deploy/static/provider/kind/deploy.yaml

kubectl patch node kind-worker --type=json -p='[{"op": "add", "path": "/metadata/labels/ingress-ready", "value": "true"}]'
}

function main {
# start-docker-desktop
 install-cluster-no-kubeproxy
 install-cilium
 install-ingress-nginx
 echo "Enable hubble UI with: cilium hubble ui &"
}

main
