package provider

import (
	"os"
	"testing"

	"github.com/devopsarr/whisparr-go/whisparr"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"whisparr": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	if v := os.Getenv("WHISPARR_URL"); v == "" {
		t.Skip("WHISPARR_URL must be set for acceptance tests")
	}

	if v := os.Getenv("WHISPARR_API_KEY"); v == "" {
		t.Skip("WHISPARR_API_KEY must be set for acceptance tests")
	}
}

func testAccAPIClient() *whisparr.APIClient {
	config := whisparr.NewConfiguration()
	config.AddDefaultHeader("X-Api-Key", os.Getenv("WHISPARR_API_KEY"))
	config.Servers[0].URL = os.Getenv("WHISPARR_URL")

	return whisparr.NewAPIClient(config)
}

const testUnauthorizedProvider = `
provider "whisparr" {
	url = "http://localhost:6969"
	api_key = "ErrorAPIKey"
  }
`
