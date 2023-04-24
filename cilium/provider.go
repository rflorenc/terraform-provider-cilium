package cilium

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mitchellh/go-homedir"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"

	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	ciliumv2alpha1 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2alpha1"
	ciliumClientset "github.com/cilium/cilium/pkg/k8s/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Register all auth providers (azure, gcp, oidc, openstack, ..).
	"k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	_ provider.Provider = &ciliumProvider{}
)

func New() provider.Provider {
	return &ciliumProvider{}
}

type ciliumProvider struct{}

type ciliumProviderModel struct {
	KubeConfig types.String `tfsdk:"kube_config"`
}

type CiliumClient struct {
	Clientset        kubernetes.Interface
	DynamicClientset dynamic.Interface
	CiliumClientset  ciliumClientset.Interface
	Config           *rest.Config
	RawConfig        clientcmdapi.Config
	restClientGetter genericclioptions.RESTClientGetter
	contextName      string
}

func NewClient(contextName, kubeconfig string) (*CiliumClient, error) {
	// Register the Cilium types in the default scheme.
	_ = ciliumv2.AddToScheme(scheme.Scheme)
	_ = ciliumv2alpha1.AddToScheme(scheme.Scheme)

	restClientGetter := genericclioptions.ConfigFlags{
		Context:    &contextName,
		KubeConfig: &kubeconfig,
	}
	rawKubeConfigLoader := restClientGetter.ToRawKubeConfigLoader()

	config, err := rawKubeConfigLoader.ClientConfig()
	if err != nil {
		return nil, err
	}

	rawConfig, err := rawKubeConfigLoader.RawConfig()
	if err != nil {
		return nil, err
	}

	ciliumClientset, err := ciliumClientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	dynamicClientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	if contextName == "" {
		contextName = rawConfig.CurrentContext
	}

	return &CiliumClient{
		CiliumClientset:  ciliumClientset,
		Clientset:        clientset,
		Config:           config,
		DynamicClientset: dynamicClientset,
		RawConfig:        rawConfig,
		restClientGetter: &restClientGetter,
		contextName:      contextName,
	}, nil
}

func (c *CiliumClient) ListCiliumNodes(ctx context.Context) (*ciliumv2.CiliumNodeList, error) {
	return c.CiliumClientset.CiliumV2().CiliumNodes().List(ctx, metav1.ListOptions{})
}

func (c *CiliumClient) ListCiliumNetworkPolicies(ctx context.Context, namespace string, opts metav1.ListOptions) (*ciliumv2.CiliumNetworkPolicyList, error) {
	return c.CiliumClientset.CiliumV2().CiliumNetworkPolicies(namespace).List(ctx, opts)
}

func (c *CiliumClient) ListCiliumClusterwideNetworkPolicies(ctx context.Context, opts metav1.ListOptions) (*ciliumv2.CiliumClusterwideNetworkPolicyList, error) {
	return c.CiliumClientset.CiliumV2().CiliumClusterwideNetworkPolicies().List(ctx, opts)
}

// Metadata should return the metadata for the provider, such as
// a type name and version data.
//
// Implementing the MetadataResponse.TypeName will populate the
// datasource.MetadataRequest.ProviderTypeName and
// resource.MetadataRequest.ProviderTypeName fields automatically.
func (hp *ciliumProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cilium"
}

// Schema should return the schema for this provider.
// The Plugin Framework uses a provider's Schema method to define the acceptable configuration attribute names and types.
func (hp *ciliumProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"kube_config": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// Configure is called at the beginning of the provider lifecycle, when
// Terraform sends to the provider the values the user specified in the
// provider configuration block. These are supplied in the
// ConfigureProviderRequest argument.
// Values from provider configuration are often used to initialise an
// API client, which should be stored on the struct implementing the
// Provider interface.
func (hp *ciliumProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring cilium client")

	// retrieve data from configuration
	var config ciliumProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.KubeConfig.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("kube_config"),
			"Unknown KubeConfig",
			"The provider cannot create the cilium API client as there is an unknown configuration value for the cilium API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the KUBE_HOST environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	kubeConfig := os.Getenv("KUBECONFIG")

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if kubeConfig == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("kube_config"),
			"Missing Kubernetes kubeConfig",
			"The provider cannot create the cilium API client as there is a missing or empty value for the cilium API password. "+
				"Set the password value in the configuration or use the KUBECONFIG environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "KUBE_CONFIG", kubeConfig)

	home, err := homedir.Dir()
	if err != nil {
		panic(err.Error())
	}
	tflog.Debug(ctx, "Creating kubernetes client")

	clientset, err := NewClient("", filepath.Join(home, ".kube", "config"))
	assertNoError(err)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Kubernetes API CiliumClient",
			"An unexpected error occurred when creating the cilium API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"cilium CiliumClient Error: "+err.Error(),
		)
		return
	}

	// Make the cilium client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = clientset
	resp.ResourceData = clientset

	tflog.Info(ctx, "Configured Kubernetes client", map[string]any{"success": true})
}

// DataSources returns a slice of functions to instantiate each DataSource
// implementation.
//
// The data source type name is determined by the DataSource implementing
// the Metadata method. All data sources must have unique names.
func (hp *ciliumProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCiliumNodeDataSource,
		NewCiliumNetworkPolicyDataSource,
		NewCiliumClusterwideNetworkPolicyDataSource,
	}
}

// Resources returns a slice of functions to instantiate each Resource
// implementation.
// The resource type name is determined by the Resource implementing
// the Metadata method. All resources must have unique names.
func (hp *ciliumProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCiliumNodeResource,
	}
}

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func assertNoErrorTF(err error, resp *datasource.ReadResponse, msg string) {
	if err != nil {
		resp.Diagnostics.AddError(
			msg,
			err.Error(),
		)
		return
	}
}
