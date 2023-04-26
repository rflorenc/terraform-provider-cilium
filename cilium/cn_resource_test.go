package cilium

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	kuberesource "k8s.io/apimachinery/pkg/api/resource"
)

func checkParsableQuantity(value string) error {
	if _, err := kuberesource.ParseQuantity(value); err != nil {
		return err
	}
	return nil
}

func TestAccCiliumNodeResource(t *testing.T) {
	nodeName := regexp.MustCompile(`^[a-z0-9]+(?:[-.]{1}[a-z0-9]+)*$`)
	zeroOrMore := regexp.MustCompile(`^[0-9]+$`)
	oneOrMore := regexp.MustCompile(`^[1-9][0-9]*$`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "cilium_ciliumnodes" "test" {
items = [
{
	# TODO
},
]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.cilium_ciliumnodes.test", "nodes.#", oneOrMore),
					resource.TestMatchResourceAttr("data.cilium_ciliumnodes.test", "nodes.0.metadata.0.labels.%", zeroOrMore),
					resource.TestCheckResourceAttrSet("data.cilium_ciliumnodes.test", "nodes.0.metadata.0.resource_version"),
					resource.TestMatchResourceAttr("data.cilium_ciliumnodes.test", "nodes.0.metadata.0.name", nodeName),
					resource.TestMatchResourceAttr("data.cilium_ciliumnodes.test", "nodes.0.spec.0.%", oneOrMore),
					resource.TestCheckResourceAttrWith("data.cilium_ciliumnodes.test", "nodes.0.status.0.capacity.cpu", checkParsableQuantity),
					resource.TestCheckResourceAttrWith("data.cilium_ciliumnodes.test", "nodes.0.status.0.capacity.memory", checkParsableQuantity),
					resource.TestCheckResourceAttrSet("data.cilium_ciliumnodes.test", "nodes.0.status.0.node_info.0.architecture"),
				),
			},
		},
	})
}

func testAccKubernetesDataSourceNodesConfig_basic() string {
	return `
data "cilium_ciliumnodes" "res" {}
`
}

func testAccKubernetesDataSourceNodesConfig_labels() string {
	return `
data "cilium_ciliumnodes" "res" {
  metadata {
    labels = {
      "kubernetes.io/os" = "linux"
    }
  }
}
`
}
