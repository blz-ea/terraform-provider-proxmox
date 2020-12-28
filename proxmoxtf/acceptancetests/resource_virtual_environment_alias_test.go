package acceptancetests

import (
	"fmt"
	"github.com/danitso/terraform-provider-proxmox/proxmoxtf"
	"github.com/danitso/terraform-provider-proxmox/proxmoxtf/acceptancetests/testutils"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strings"
	"testing"
)

// Verifies that an alias can be created and updated
func TestAccResourceVirtualEnvironmentAlias_CreateAndUpdate(t *testing.T) {
	aliasNameFirst := testutils.GenerateResourceName()
	aliasNameSecond := testutils.GenerateResourceName()

	tfNode := "proxmox_virtual_environment_cluster_alias.alias"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckAliasDestroyed,
		Steps: []resource.TestStep{
			// Create alias
			{
				Config: testutils.HclAliasResource(aliasNameFirst, "192.168.0.0/23", "alias-comment"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", aliasNameFirst),
					resource.TestCheckResourceAttr(tfNode, "cidr", "192.168.0.0/23"),
					resource.TestCheckResourceAttr(tfNode, "comment", "alias-comment"),
					testutils.CheckAliasExists(aliasNameFirst),
				),
			},
			// Update alias
			{
				Config: testutils.HclAliasResource(aliasNameSecond, "192.168.0.1", "alias-comment"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", aliasNameSecond),
					resource.TestCheckResourceAttr(tfNode, "cidr", "192.168.0.1"),
					resource.TestCheckResourceAttr(tfNode, "comment", "alias-comment"),
					testutils.CheckAliasExists(aliasNameSecond),
				),
			},
		},
	})
}

// CheckAliasDestroyed verifies that all aliases referenced in the state
// are destroyed. This will be invoked *after* terraform destroys
// the resource but *before* the state is wiped clean
func CheckAliasDestroyed(s *terraform.State) error {
	config := testutils.GetProvider().Meta().(proxmoxtf.ProviderConfiguration)

	conn, err := config.GetVEClient()

	if err != nil {
		return err
	}

	// loop through the resource state
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "proxmox_virtual_environment_cluster_alias" {
			continue
		}

		response, err := conn.GetAlias(rs.Primary.ID)

		if err == nil {
			if response.Name != "" && response.Name == rs.Primary.ID {
				return fmt.Errorf("Alias with Name=`%s` should not exist", rs.Primary.ID)
			}

			return nil
		}

		// If the error is not 400 (which identifies if role was found or not)
		if !strings.Contains(err.Error(), "verification failed") {
			return err
		}
	}

	return nil
}
