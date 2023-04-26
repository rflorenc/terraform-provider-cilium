package cilium

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCoffeesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "cilium_ciliumnodes" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of ciliumnodes returned
					resource.TestCheckResourceAttr("data.cilium_ciliumnodes.test", "ciliumnode.#", "3"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("data.cilium_ciliumnodes.test", "id", "placeholder"),
				),
			},
		},
	})
}
