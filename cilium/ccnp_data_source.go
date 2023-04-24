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
	_ datasource.DataSource              = &ciliumClusterwideNetworkPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &ciliumClusterwideNetworkPolicyDataSource{}
)

// NewCiliumClusterwideNetworkPolicyDataSource is a helper function to simplify the provider implementation.
func NewCiliumClusterwideNetworkPolicyDataSource() datasource.DataSource {
	return &ciliumClusterwideNetworkPolicyDataSource{}
}

// ciliumClusterwideNetworkPolicyDataSource is the data source implementation.
type ciliumClusterwideNetworkPolicyDataSource struct {
	client *CiliumClient
}

// Metadata returns the data source type name.
func (d *ciliumClusterwideNetworkPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ciliumClusterwideNetworkPolicies"
}

// Schema defines the schema for the data source.
func (d *ciliumClusterwideNetworkPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ciliumclusterwidenetworkpolicies": schema.ListNestedAttribute{
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
func (d *ciliumClusterwideNetworkPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ccnpl, err := d.client.ListCiliumClusterwideNetworkPolicies(ctx, metav1.ListOptions{})
	assertNoError(err)

	tflog.Debug(ctx, "DBG CCNPL", map[string]interface{}{
		"cnpnl result": ccnpl,
	})

	tflog.Debug(ctx, "DBG CCNPL items", map[string]interface{}{
		"ccnpl.Items: ": ccnpl.Items,
	})
}

// Configure enables provider-level data or clients to be set in the
// provider-defined DataSource type. It is separately executed for each
// ReadDataSource RPC.
func (d *ciliumClusterwideNetworkPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*CiliumClient)
}

// coffeesDataSourceModel maps the data source schema data.
type ciliumClusterwideNetworkPolicysDataSourceModel struct {
	ciliumNetworkPolicys []ciliumClusterwideNetworkPolicyModel `tfsdk:"ciliumclusterwidenetworkpolicies"`
}

// ciliumClusterwideNetworkPolicyModel maps ciliumNetworkPolicys schema data.
type ciliumClusterwideNetworkPolicyModel struct {
	ApiVersion types.String `tfsdk:"apiVersion"`
	Kind       types.String `tfsdk:"kind"`
}
