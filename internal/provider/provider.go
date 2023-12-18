// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"
	"os"
	"time"

	"experimental-terraform-redirect-store/api"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure RedirectStoreProvider satisfies various provider interfaces.
var _ provider.Provider = &RedirectStoreProvider{}

// RedirectStoreProvider defines the provider implementation.
type RedirectStoreProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// RedirectStoreProviderModel describes the provider data model.
type RedirectStoreProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *RedirectStoreProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "redirect-store"
	resp.Version = p.version
}

func (p *RedirectStoreProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Redirect Store",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "API endpoint",
				Optional:            true,
			},
		},
	}
}

func (p *RedirectStoreProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config RedirectStoreProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown RedirectStore API Endpoint",
			"The provider cannot create the RedirectStore API client as there is an unknown configuration value for the RedirectStore API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the REDIRECT_STORE_ENDPOINT environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("REDIRECT_STORE_ENDPOINT")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing ResirectStore API Endpoint",
			"The provider cannot create the ResirectStore API client as there is a missing or empty value for the RedirectStore API host. "+
				"Set the endpoint value in the configuration or use the REDIRECT_STORE_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client := api.NewClientImpl(endpoint, &http.Client{
		Timeout: 3 * time.Second,
	})
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured RedirectStore client", map[string]any{"endpoint": endpoint})
}

func (p *RedirectStoreProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRecordResource,
	}
}

func (p *RedirectStoreProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRecordsDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RedirectStoreProvider{
			version: version,
		}
	}
}
