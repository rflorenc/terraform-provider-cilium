package cilium

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ciliumNodeDataSource{}
	_ datasource.DataSourceWithConfigure = &ciliumNodeDataSource{}
)

// NewCiliumNodeDataSource is a helper function to simplify the provider implementation.
func NewCiliumNodeDataSource() datasource.DataSource {
	return &ciliumNodeDataSource{}
}

// ciliumNodeDataSource is the data source implementation.
type ciliumNodeDataSource struct {
	client *CiliumClient
}

// Metadata returns the data source type name.
func (d *ciliumNodeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ciliumnodes"
}

// Schema defines the schema for the data source.
func (d *ciliumNodeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ciliumnodes": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"apiversion": schema.StringAttribute{
							Computed: true,
						},
						"kind": schema.StringAttribute{
							Computed: true,
						},
						"metadata": schema.ObjectAttribute{
							Computed: true,
						},
						"spec": schema.ObjectAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *ciliumNodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ciliumNodeDataSourceModel

	cnl, err := d.client.ListCiliumNodes(ctx)
	assertNoErrorTF(err, resp, "Unable to ListCiliumNodes")

	for _, cn := range cnl.Items {
		ciliumNodeState := ciliumNodeModel{
			ApiVersion: types.StringValue(cn.APIVersion),
			Kind:       types.StringValue(cn.Kind),
			Metadata: &metadataModel{
				annotations: cn.Annotations,
				labels:      cn.Labels,
				name:        types.StringValue(cn.Name),
				namespace:   types.StringValue(cn.Namespace),
			},
			Spec: &specModel{
				instanceID: types.StringValue(cn.Spec.InstanceID),
			},
		}
		tflog.Debug(ctx, "DBG state.CiliumNodes:::", map[string]interface{}{
			"tf state:: ":          state.CiliumNodes,
			"tf ciliumNodeState: ": ciliumNodeState,
		})
		state.CiliumNodes = append(state.CiliumNodes, ciliumNodeState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure enables provider-level data or clients to be set in the
// provider-defined DataSource type. It is separately executed for each
// ReadDataSource RPC.
func (d *ciliumNodeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*CiliumClient)
}

// coffeesDataSourceModel maps the data source schema data.
type ciliumNodeDataSourceModel struct {
	CiliumNodes []ciliumNodeModel `tfsdk:"ciliumnodes"`
}

// ciliumNodeModel maps ciliumnodes schema data.
type ciliumNodeModel struct {
	ApiVersion types.String   `tfsdk:"apiversion"`
	Kind       types.String   `tfsdk:"kind"`
	Metadata   *metadataModel `tfsdk:"metadata"`
	Spec       *specModel     `tfsdk:"spec"`
}

type metadataModel struct {
	annotations map[string]string `tfsdk:"annotations"`
	labels      map[string]string `tfsdk:"labels"`
	name        types.String      `tfsdk:"name"`
	namespace   types.String      `tfsdk:"namespace"`
}

type specModel struct {
	instanceID types.String `tfsdk:"instance-id"`
}
