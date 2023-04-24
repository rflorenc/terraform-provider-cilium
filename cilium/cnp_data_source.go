package cilium

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ciliumNetworkPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &ciliumNetworkPolicyDataSource{}
)

// NewCiliumNetworkPolicyDataSource is a helper function to simplify the provider implementation.
func NewCiliumNetworkPolicyDataSource() datasource.DataSource {
	return &ciliumNetworkPolicyDataSource{}
}

// ciliumNetworkPolicyDataSource is the data source implementation.
type ciliumNetworkPolicyDataSource struct {
	client *CiliumClient
}

// Metadata returns the data source type name.
func (d *ciliumNetworkPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ciliumNetworkPolicies"
}

// Schema defines the schema for the data source.
func (d *ciliumNetworkPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ciliumnetworkpolicies": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"apiversion": schema.StringAttribute{
							Computed: true,
						},
						"kind": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *ciliumNetworkPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	cnpl, err := d.client.ListCiliumNetworkPolicies(ctx, "", metav1.ListOptions{})
	assertNoError(err)

	tflog.Debug(ctx, "DBG CNPL", map[string]interface{}{
		"cn result": cnpl,
	})

	tflog.Debug(ctx, "DBG CNPL items", map[string]interface{}{
		"cnp.Name: ": cnpl.Items,
	})
}

// Configure enables provider-level data or clients to be set in the
// provider-defined DataSource type. It is separately executed for each
// ReadDataSource RPC.
func (d *ciliumNetworkPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*CiliumClient)
}

// coffeesDataSourceModel maps the data source schema data.
type ciliumNetworkPoliciesDataSourceModel struct {
	ciliumNetworkPolicies []ciliumNetworkPolicyModel `tfsdk:"ciliumnetworkpolicies"`
}

// ciliumNetworkPolicyModel maps ciliumNetworkPolicys schema data.
type ciliumNetworkPolicyModel struct {
	ApiVersion types.String `tfsdk:"apiVersion"`
	Kind       types.String `tfsdk:"kind"`
}
