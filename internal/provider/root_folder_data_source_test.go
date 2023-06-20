package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/devopsarr/whisparr-go/whisparr"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRootFolderDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccRootFolderDataSourceConfig("/error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccRootFolderDataSourceConfig("/error"),
				ExpectError: regexp.MustCompile("Unable to find root_folder"),
			},
			// Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccRootFolderDataSourceConfig("/config"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.whisparr_root_folder.test", "id"),
					resource.TestCheckResourceAttr("data.whisparr_root_folder.test", "path", "/config")),
			},
		},
	})
}

func testAccRootFolderDataSourceConfig(path string) string {
	return fmt.Sprintf(`
	data "whisparr_root_folder" "test" {
  			path = "%s"
		}
	`, path)
}

func rootFolderDSInit() {
	// ensure a /config root path is configured
	client := testAccAPIClient()
	folder := whisparr.NewRootFolderResource()
	folder.SetPath("/config")
	_, _, _ = client.RootFolderApi.CreateRootFolder(context.TODO()).RootFolderResource(*folder).Execute()
}
