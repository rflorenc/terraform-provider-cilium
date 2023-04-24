package cilium

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &ciliumNodeResource{}
	//_ resource.ResourceWithConfigure   = &orderResource{}
	//_ resource.ResourceWithImportState = &orderResource{}
)

// NewCiliumNodeResource is a helper function to simplify the provider implementation.
func NewCiliumNodeResource() resource.Resource {
	return &ciliumNodeResource{}
}

// orderResourceModel maps the resource schema data.
type ciliumNodeResourceModel struct {
	ID types.String `tfsdk:"id"`
}

// ciliumNodeResource is the resource implementation.
type ciliumNodeResource struct {
	client *CiliumClient
}

// Configure enables provider-level data or clients to be set in the
// provider-defined DataSource type. It is separately executed for each
// ReadDataSource RPC.
func (r *ciliumNodeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*CiliumClient)
}

// Metadata returns the resource type name.
func (r *ciliumNodeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_order"
}

// Schema defines the schema for the resource.
func (r *ciliumNodeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

// Create creates the resource and sets the initial Terraform state.
func (r *ciliumNodeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ciliumNodeResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *ciliumNodeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ciliumNodeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
// https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-resource-update#implement-update-functionality
func (r *ciliumNodeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan ciliumNodeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-resource-delete#implement-delete-functionality
func (r *ciliumNodeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state

}

func (r *ciliumNodeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
