package testutils

import (
	"fmt"
	"github.com/blz-ea/terraform-provider-proxmox/proxmox"
	"github.com/blz-ea/terraform-provider-proxmox/proxmoxtf"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// CheckAliasExists Given the name of alias, this will return a function that will check
// whether or not an alias
// - (1) exists in the state
// - (2) exist in Proxmox VE
// - (3) has the correct name
func CheckAliasExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources["proxmox_virtual_environment_cluster_alias.alias"]
		if !ok {
			return fmt.Errorf("Did not find the alias in the TF state")
		}

		clients := GetProvider().Meta().(proxmoxtf.ProviderConfiguration)
		id := res.Primary.ID
		alias, err := readAlias(clients, id)

		if err != nil {
			return fmt.Errorf("Alias with Name=%s cannot be found. Error %v", id, err)
		}

		if alias.Name != expectedName {
			return fmt.Errorf("Alias with Name=%s has Name=%s, but expected Name=%s", id, alias.Name, expectedName)
		}

		return nil
	}
}

// readAlias is a helper function that reads an alias based on a given name
func readAlias(clients proxmoxtf.ProviderConfiguration, identifier string) (*proxmox.VirtualEnvironmentClusterAliasGetResponseData, error) {
	conn, err := clients.GetVEClient()

	if err != nil {
		return nil, err
	}

	response, err := conn.GetAlias(identifier)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// HclAliasResource HCL describing of a PVE alias resource
func HclAliasResource(name string, cidr string, comment string) string {

	if name == "" {
		panic("Parameter: `name` cannot be empty")
	}

	if cidr == "" {
		panic("Parameter: `cidr` cannot be empty")
	}

	return fmt.Sprintf(`
resource "proxmox_virtual_environment_cluster_alias" "alias" {
	name    = "%[1]s"
	cidr    = "%[2]s"
	comment = "%[3]s"
}
`, name, cidr, comment)
}
