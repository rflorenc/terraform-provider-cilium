# Cilium Terraform Provider


Use this provider as a bridge between Terraform and the cilium kubernetes CRDs.

## Build provider

Run the following command to build the provider

```shell
$ go build -o terraform-provider-cilium
```

## Test sample configuration

First, build and install the provider.

```shell
$ make install
```

Then, navigate to the `examples` directory.

```shell
$ cd examples
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ export KUBECONFIG="~/.kube/config"
$ terraform init && terraform apply
```
