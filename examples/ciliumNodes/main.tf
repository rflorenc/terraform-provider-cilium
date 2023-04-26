terraform {
  required_providers {
    cilium = {
      source = "hashicorp.com/edu/cilium"
    }
  }
}

provider "cilium" {
  kube_config = "~/.kube/config"
}

data "cilium_ciliumnodes" "res" {}

output "edu_ciliumnodes" {
  value = data.cilium_ciliumnodes.res
}
