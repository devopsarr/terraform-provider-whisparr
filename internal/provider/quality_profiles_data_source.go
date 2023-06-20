package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/terraform-provider-whisparr/internal/helpers"
	"github.com/devopsarr/whisparr-go/whisparr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const qualityProfilesDataSourceName = "quality_profiles"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &QualityProfilesDataSource{}

func NewQualityProfilesDataSource() datasource.DataSource {
	return &QualityProfilesDataSource{}
}

// QualityProfilesDataSource defines the qyality profiles implementation.
type QualityProfilesDataSource struct {
	client *whisparr.APIClient
}

// QualityProfiles describes the qyality profiles data model.
type QualityProfiles struct {
	QualityProfiles types.Set    `tfsdk:"quality_profiles"`
	ID              types.String `tfsdk:"id"`
}

func (d *QualityProfilesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + qualityProfilesDataSourceName
}

func (d *QualityProfilesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the quality server.
		MarkdownDescription: "<!-- subcategory:Profiles -->List all available [Quality Profiles](../resources/quality_profile).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"quality_profiles": schema.SetNestedAttribute{
				MarkdownDescription: "Quality Profile list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "Quality Profile ID.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Quality Profile Name.",
							Computed:            true,
						},
						"upgrade_allowed": schema.BoolAttribute{
							MarkdownDescription: "Upgrade allowed flag.",
							Computed:            true,
						},
						"cutoff": schema.Int64Attribute{
							MarkdownDescription: "Quality ID to which cutoff.",
							Computed:            true,
						},
						"cutoff_format_score": schema.Int64Attribute{
							MarkdownDescription: "Cutoff format score.",
							Computed:            true,
						},
						"min_format_score": schema.Int64Attribute{
							MarkdownDescription: "Min format score.",
							Computed:            true,
						},
						"language": schema.SingleNestedAttribute{
							MarkdownDescription: "Language.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"id": schema.Int64Attribute{
									MarkdownDescription: "ID.",
									Computed:            true,
								},
								"name": schema.StringAttribute{
									MarkdownDescription: "Name.",
									Computed:            true,
								},
							},
						},
						"quality_groups": schema.SetNestedAttribute{
							MarkdownDescription: "Quality groups.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.Int64Attribute{
										MarkdownDescription: "Quality group ID.",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										MarkdownDescription: "Quality group name.",
										Computed:            true,
									},
									"qualities": schema.SetNestedAttribute{
										MarkdownDescription: "Qualities in group.",
										Required:            true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"id": schema.Int64Attribute{
													MarkdownDescription: "Quality ID.",
													Computed:            true,
												},
												"resolution": schema.Int64Attribute{
													MarkdownDescription: "Resolution.",
													Computed:            true,
												},
												"name": schema.StringAttribute{
													MarkdownDescription: "Quality name.",
													Computed:            true,
												},
												"source": schema.StringAttribute{
													MarkdownDescription: "Source.",
													Computed:            true,
												},
											},
										},
									},
								},
							},
						},
						"format_items": schema.SetNestedAttribute{
							MarkdownDescription: "Format items.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"format": schema.Int64Attribute{
										MarkdownDescription: "Format.",
										Computed:            true,
									},
									"score": schema.Int64Attribute{
										MarkdownDescription: "Score.",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										MarkdownDescription: "Name.",
										Computed:            true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *QualityProfilesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *QualityProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *QualityProfiles

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get qualityprofiles current value
	response, _, err := d.client.QualityProfileApi.ListQualityProfile(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.List, qualityProfilesDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+qualityProfilesDataSourceName)
	// Map response body to resource schema attribute
	profiles := *writeQualitiyprofiles(ctx, response)
	tfsdk.ValueFrom(ctx, profiles, data.QualityProfiles.Type(ctx), &data.QualityProfiles)

	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func writeQualitiyprofiles(ctx context.Context, qualities []*whisparr.QualityProfileResource) *[]QualityProfile {
	output := make([]QualityProfile, len(qualities))
	for i, p := range qualities {
		output[i].write(ctx, p, &diag.Diagnostics{})
	}

	return &output
}
