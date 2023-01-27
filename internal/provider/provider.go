package provider

import (
	"context"
	"os"

	"github.com/devopsarr/whisparr-go/whisparr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// needed for tf debug mode
// var stderr = os.Stderr

// Ensure provider defined types fully satisfy framework interfaces.
var _ provider.Provider = &WhisparrProvider{}

// ScaffoldingProvider defines the provider implementation.
type WhisparrProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Whisparr describes the provider data model.
type Whisparr struct {
	APIKey types.String `tfsdk:"api_key"`
	URL    types.String `tfsdk:"url"`
}

func (p *WhisparrProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "whisparr"
	resp.Version = p.version
}

func (p *WhisparrProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Whisparr provider is used to interact with any [Whisparr](https://whisparr.video/) installation. You must configure the provider with the proper credentials before you can use it. Use the left navigation to read about the available resources.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key for Whisparr authentication. Can be specified via the `WHISPARR_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Full Whisparr URL with protocol and port (e.g. `https://test.whisparr.tv:6969`). You should **NOT** supply any path (`/api`), the SDK will use the appropriate paths. Can be specified via the `WHISPARR_URL` environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *WhisparrProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data Whisparr

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// User must provide URL to the provider
	if data.URL.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as url",
		)

		return
	}

	var url string
	if data.URL.IsNull() {
		url = os.Getenv("WHISPARR_URL")
	} else {
		url = data.URL.ValueString()
	}

	if url == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find URL",
			"URL cannot be an empty string",
		)

		return
	}

	// User must provide API key to the provider
	if data.APIKey.IsUnknown() {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as api_key",
		)

		return
	}

	var key string
	if data.APIKey.IsNull() {
		key = os.Getenv("WHISPARR_API_KEY")
	} else {
		key = data.APIKey.ValueString()
	}

	if key == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find API key",
			"API key cannot be an empty string",
		)

		return
	}

	// Configuring client. API Key management could be changed once new options avail in sdk.
	config := whisparr.NewConfiguration()
	config.AddDefaultHeader("X-Api-Key", key)
	config.Servers[0].URL = url
	client := whisparr.NewAPIClient(config)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *WhisparrProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Download Clients
		NewDownloadClientConfigResource,
		NewDownloadClientResource,
		NewDownloadClientTransmissionResource,
		NewDownloadClientAria2Resource,
		NewDownloadClientDelugeResource,
		NewDownloadClientFloodResource,
		NewDownloadClientHadoukenResource,
		NewDownloadClientNzbgetResource,
		NewDownloadClientNzbvortexResource,
		NewDownloadClientPneumaticResource,
		NewDownloadClientQbittorrentResource,
		NewDownloadClientRtorrentResource,
		NewDownloadClientSabnzbdResource,
		NewDownloadClientTorrentBlackholeResource,
		NewDownloadClientTorrentDownloadStationResource,
		NewDownloadClientUsenetBlackholeResource,
		NewDownloadClientUsenetDownloadStationResource,
		NewDownloadClientUtorrentResource,
		NewDownloadClientVuzeResource,
		NewRemotePathMappingResource,

		// Indexers
		NewIndexerConfigResource,
		NewIndexerResource,
		NewIndexerFilelistResource,
		NewIndexerIptorrentsResource,
		NewIndexerHdbitsResource,
		NewIndexerNewznabResource,
		NewIndexerNyaaResource,
		NewIndexerOmgwtfnzbsResource,
		NewIndexerRarbgResource,
		NewIndexerTorrentPotatoResource,
		NewIndexerTorrentRssResource,
		NewIndexerTorznabResource,
		NewRestrictionResource,

		// Import Lists
		NewImportListConfigResource,
		NewImportListResource,
		NewImportListCouchPotatoResource,
		NewImportListCustomResource,
		NewImportListIMDBResource,
		NewImportListPlexResource,
		NewImportListRSSResource,
		NewImportListStevenluResource,
		NewImportListStevenlu2Resource,
		NewImportListTMDBCollectionResource,
		NewImportListTMDBCompanyResource,
		NewImportListTMDBKeywordResource,
		NewImportListTMDBListResource,
		NewImportListTMDBPersonResource,
		NewImportListTMDBPopularResource,
		NewImportListTMDBUserResource,
		NewImportListTraktListResource,
		NewImportListTraktPopularResource,
		NewImportListTraktUserResource,
		NewImportListWhisparrResource,
		NewImportListExclusionResource,

		// Media Management
		NewMediaManagementResource,
		NewNamingResource,
		NewRootFolderResource,

		// Metadata
		NewMetadataResource,
		NewMetadataEmbyResource,
		NewMetadataKodiResource,
		NewMetadataRoksboxResource,
		NewMetadataWdtvResource,
		NewMetadataConfigResource,

		// Notifications
		NewNotificationResource,
		NewNotificationBoxcarResource,
		NewNotificationCustomScriptResource,
		NewNotificationDiscordResource,
		NewNotificationEmailResource,
		NewNotificationEmbyResource,
		NewNotificationGotifyResource,
		NewNotificationJoinResource,
		NewNotificationKodiResource,
		NewNotificationMailgunResource,
		NewNotificationNotifiarrResource,
		NewNotificationPlexResource,
		NewNotificationProwlResource,
		NewNotificationPushbulletResource,
		NewNotificationPushoverResource,
		NewNotificationSendgridResource,
		NewNotificationSimplepushResource,
		NewNotificationSlackResource,
		NewNotificationSynologyResource,
		NewNotificationTelegramResource,
		NewNotificationTraktResource,
		NewNotificationTwitterResource,
		NewNotificationWebhookResource,

		// Profiles
		NewCustomFormatResource,
		NewDelayProfileResource,
		NewQualityProfileResource,

		// Tags
		NewTagResource,
	}
}

func (p *WhisparrProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Download Clients
		NewDownloadClientConfigDataSource,
		NewDownloadClientDataSource,
		NewDownloadClientsDataSource,
		NewRemotePathMappingDataSource,
		NewRemotePathMappingsDataSource,

		// Indexers
		NewIndexerConfigDataSource,
		NewIndexerDataSource,
		NewIndexersDataSource,
		NewRestrictionDataSource,
		NewRestrictionsDataSource,

		// Import Lists
		NewImportListConfigDataSource,
		NewImportListDataSource,
		NewImportListsDataSource,
		NewImportListExclusionDataSource,
		NewImportListExclusionsDataSource,

		// Media Management
		NewMediaManagementDataSource,
		NewNamingDataSource,
		NewRootFolderDataSource,
		NewRootFoldersDataSource,

		// Metadata
		NewMetadataConfigDataSource,
		NewMetadataDataSource,
		NewMetadataConsumersDataSource,

		// Notifications
		NewNotificationDataSource,
		NewNotificationsDataSource,

		// Profiles
		NewCustomFormatDataSource,
		NewCustomFormatsDataSource,
		NewDelayProfileDataSource,
		NewDelayProfilesDataSource,
		NewQualityProfileDataSource,
		NewQualityProfilesDataSource,

		// System Status
		NewSystemStatusDataSource,

		// Tags
		NewTagDataSource,
		NewTagsDataSource,
	}
}

// New returns the provider with a specific version.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &WhisparrProvider{
			version: version,
		}
	}
}
