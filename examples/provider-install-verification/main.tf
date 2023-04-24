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

data "cilium_ciliumNode" "example" {}
