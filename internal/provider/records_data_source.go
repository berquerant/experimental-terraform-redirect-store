package provider

import (
	"context"
	"errors"
	"experimental-terraform-redirect-store/api"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &recordsDataSource{}
	_ datasource.DataSourceWithConfigure = &recordsDataSource{}
)

func NewRecordsDataSource() datasource.DataSource {
	return &recordsDataSource{}
}

type recordsDataSource struct {
	client api.Client
}

type recordsDataSourceModel struct {
	Records []recordsModel `tfsdk:"records"`
}

type recordsModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	To   types.String `tfsdk:"to"`
}

func (d *recordsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_records"
}

func (d *recordsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch the list of records.",
		Attributes: map[string]schema.Attribute{
			"records": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Placeholder identifier attribute.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Record name.",
							Computed:    true,
						},
						"to": schema.StringAttribute{
							Description: "Record redirect-to.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *recordsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state recordsDataSourceModel

	records, err := d.client.Scan(ctx)
	switch {
	case errors.Is(err, api.ErrNotFound):
		// not found but ok
	case err != nil:
		resp.Diagnostics.AddError(
			"Unable to Read RedirectStore Records",
			err.Error(),
		)
		return
	default:
		for _, record := range records {
			state.Records = append(state.Records, recordsModel{
				ID:   types.StringValue(record.Name),
				Name: types.StringValue(record.Name),
				To:   types.StringValue(record.To),
			})
		}
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *recordsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
