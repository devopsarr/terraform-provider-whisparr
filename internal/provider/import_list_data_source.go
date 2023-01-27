package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/terraform-provider-whisparr/internal/helpers"
	"github.com/devopsarr/whisparr-go/whisparr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const importListDataSourceName = "import_list"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ImportListDataSource{}

func NewImportListDataSource() datasource.DataSource {
	return &ImportListDataSource{}
}

// ImportListDataSource defines the import_list implementation.
type ImportListDataSource struct {
	client *whisparr.APIClient
}

func (d *ImportListDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListDataSourceName
}

func (d *ImportListDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Import Lists -->Single [Import List](../resources/import_list).",
		Attributes: map[string]schema.Attribute{
			"enable_auto": schema.BoolAttribute{
				MarkdownDescription: "Enable automatic add flag.",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Enabled flag.",
				Computed:            true,
			},
			"search_on_add": schema.BoolAttribute{
				MarkdownDescription: "Search on add flag.",
				Computed:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality profile ID.",
				Computed:            true,
			},
			"list_order": schema.Int64Attribute{
				MarkdownDescription: "List order.",
				Computed:            true,
			},
			"root_folder_path": schema.StringAttribute{
				MarkdownDescription: "Root folder path.",
				Computed:            true,
			},
			"should_monitor": schema.BoolAttribute{
				MarkdownDescription: "Should monitor.",
				Computed:            true,
			},
			"minimum_availability": schema.StringAttribute{
				MarkdownDescription: "Minimum availability.",
				Computed:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "ImportList implementation name.",
				Computed:            true,
			},
			"config_contract": schema.StringAttribute{
				MarkdownDescription: "ImportList configuration template.",
				Computed:            true,
			},
			"list_type": schema.StringAttribute{
				MarkdownDescription: "List type.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Import List name.",
				Required:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Import List ID.",
				Computed:            true,
			},
			// Field values
			"only_active": schema.BoolAttribute{
				MarkdownDescription: "Only active.",
				Computed:            true,
			},
			"cast": schema.BoolAttribute{
				MarkdownDescription: "Include cast.",
				Computed:            true,
			},
			"cast_director": schema.BoolAttribute{
				MarkdownDescription: "Include cast director.",
				Computed:            true,
			},
			"cast_producer": schema.BoolAttribute{
				MarkdownDescription: "Include cast producer.",
				Computed:            true,
			},
			"cast_sound": schema.BoolAttribute{
				MarkdownDescription: "Include cast sound.",
				Computed:            true,
			},
			"cast_writing": schema.BoolAttribute{
				MarkdownDescription: "Include cast writing.",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port.",
				Computed:            true,
			},
			"source": schema.Int64Attribute{
				MarkdownDescription: "Source.",
				Computed:            true,
			},
			"min_score": schema.Int64Attribute{
				MarkdownDescription: "Min score.",
				Computed:            true,
			},
			"tmdb_list_type": schema.Int64Attribute{
				MarkdownDescription: "TMDB list type.",
				Computed:            true,
			},
			"user_list_type": schema.Int64Attribute{
				MarkdownDescription: "User list type.",
				Computed:            true,
			},
			"limit": schema.Int64Attribute{
				MarkdownDescription: "limit.",
				Computed:            true,
			},
			"trakt_list_type": schema.Int64Attribute{
				MarkdownDescription: "Trakt list type.",
				Computed:            true,
			},
			"language_code": schema.Int64Attribute{
				MarkdownDescription: "Language code.",
				Computed:            true,
			},
			"listname": schema.StringAttribute{
				MarkdownDescription: "List name.",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username.",
				Computed:            true,
			},
			"auth_user": schema.StringAttribute{
				MarkdownDescription: "Auth user.",
				Computed:            true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Access token.",
				Computed:            true,
				Sensitive:           true,
			},
			"refresh_token": schema.StringAttribute{
				MarkdownDescription: "Refresh token.",
				Computed:            true,
				Sensitive:           true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key.",
				Computed:            true,
				Sensitive:           true,
			},
			"company_id": schema.StringAttribute{
				MarkdownDescription: "Company ID.",
				Computed:            true,
			},
			"collection_id": schema.StringAttribute{
				MarkdownDescription: "Collection ID.",
				Computed:            true,
			},
			"keyword_id": schema.StringAttribute{
				MarkdownDescription: "Keyword ID.",
				Computed:            true,
			},
			"list_id": schema.StringAttribute{
				MarkdownDescription: "List ID.",
				Computed:            true,
			},
			"person_id": schema.StringAttribute{
				MarkdownDescription: "Person ID.",
				Computed:            true,
			},
			"account_id": schema.StringAttribute{
				MarkdownDescription: "Account ID.",
				Computed:            true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Computed:            true,
			},
			"url_base": schema.StringAttribute{
				MarkdownDescription: "Base URL.",
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL.",
				Computed:            true,
			},
			"link": schema.StringAttribute{
				MarkdownDescription: "Link.",
				Computed:            true,
			},
			"expires": schema.StringAttribute{
				MarkdownDescription: "Expires.",
				Computed:            true,
			},
			"trakt_additional_parameters": schema.StringAttribute{
				MarkdownDescription: "Trakt additional parameters.",
				Computed:            true,
			},
			"certification": schema.StringAttribute{
				MarkdownDescription: "Certification.",
				Computed:            true,
			},
			"genres": schema.StringAttribute{
				MarkdownDescription: "Genres.",
				Computed:            true,
			},
			"years": schema.StringAttribute{
				MarkdownDescription: "Years.",
				Computed:            true,
			},
			"rating": schema.StringAttribute{
				MarkdownDescription: "Rating.",
				Computed:            true,
			},
			"min_vote_average": schema.StringAttribute{
				MarkdownDescription: "Min vote average.",
				Computed:            true,
			},
			"min_votes": schema.StringAttribute{
				MarkdownDescription: "Min votes.",
				Computed:            true,
			},
			"tmdb_certification": schema.StringAttribute{
				MarkdownDescription: "Certification.",
				Computed:            true,
			},
			"include_genre_ids": schema.StringAttribute{
				MarkdownDescription: "Include genre IDs.",
				Computed:            true,
			},
			"exclude_genre_ids": schema.StringAttribute{
				MarkdownDescription: "Exclude genre IDs.",
				Computed:            true,
			},
			"profile_ids": schema.SetAttribute{
				MarkdownDescription: "Profile IDs.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"tag_ids": schema.SetAttribute{
				MarkdownDescription: "Tag IDs.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
		},
	}
}

func (d *ImportListDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *ImportListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ImportList

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get importList current value
	response, _, err := d.client.ImportListApi.ListImportList(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListDataSourceName, err))

		return
	}

	importList, err := findImportList(data.Name.ValueString(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", importListDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListDataSourceName)
	data.write(ctx, importList)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findImportList(name string, importLists []*whisparr.ImportListResource) (*whisparr.ImportListResource, error) {
	for _, i := range importLists {
		if i.GetName() == name {
			return i, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(importListDataSourceName, "name", name)
}
