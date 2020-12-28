package acceptancetests

import (
	"fmt"
	"github.com/blz-ea/terraform-provider-proxmox/proxmox"
	"github.com/blz-ea/terraform-provider-proxmox/proxmoxtf"
	"github.com/blz-ea/terraform-provider-proxmox/proxmoxtf/acceptancetests/testutils"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strings"
	"testing"
)

// Verifies that an IP set can be created and updated
func TestAccResourceVirtualEnvironmentIPSet_CreateAndUpdate(t *testing.T) {
	IPSetNameFirst := testutils.GenerateResourceName()
	IPSetCIDREmpty := proxmox.VirtualEnvironmentClusterIPSetContent{}
	noMatch := proxmox.CustomBool(true)
	IPSetCIDR := proxmox.VirtualEnvironmentClusterIPSetContent{
		{CIDR: "192.168.88.1", Comment: "ipset-cidr-comment"},
		{CIDR: "192.168.88.2", Comment: "ipset-cidr-comment", NoMatch: &noMatch},
	}

	IPSetCIDRWithAlias := proxmox.VirtualEnvironmentClusterIPSetContent{
		{CIDR: "192.168.88.4", Comment: "ipset-cidr-comment"},
		// Alias resource is created by HclIPSetWithAliasResource
		{CIDR: "test-alias", Comment: "ipset-cidr-comment", NoMatch: &noMatch},
	}

	IPSetFirst := proxmox.VirtualEnvironmentClusterIPSetContent{IPSetCIDR[0]}
	IPSetSecond := proxmox.VirtualEnvironmentClusterIPSetContent{IPSetCIDR[1]}

	IPSetNameSecond := testutils.GenerateResourceName()
	tfNode := "proxmox_virtual_environment_cluster_ipset.ipset"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckIPSetDestroyed,
		Steps: []resource.TestStep{
			// Create empty IPSet
			{
				Config: testutils.HclIPSetResource(IPSetNameFirst, "ipset-comment", IPSetCIDREmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", IPSetNameFirst),
					resource.TestCheckResourceAttr(tfNode, "comment", "ipset-comment"),
					testutils.CheckIPSetExists(IPSetCIDREmpty),
				),
			},
			// Update IPSet with CIDR
			{
				Config: testutils.HclIPSetResource(IPSetNameFirst, "ipset-comment", IPSetFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", IPSetNameFirst),
					resource.TestCheckResourceAttr(tfNode, "comment", "ipset-comment"),
					testutils.CheckIPSetExists(IPSetFirst),
				),
			},
			// Update IPSet's name and comment
			{
				Config: testutils.HclIPSetResource(IPSetNameSecond, "ipset-comment-updated", IPSetFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", IPSetNameSecond),
					resource.TestCheckResourceAttr(tfNode, "comment", "ipset-comment-updated"),
					testutils.CheckIPSetExists(IPSetFirst),
				),
			},
			// Update IPSet's CIDR
			{
				Config: testutils.HclIPSetResource(IPSetNameSecond, "ipset-comment-updated", IPSetSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", IPSetNameSecond),
					resource.TestCheckResourceAttr(tfNode, "comment", "ipset-comment-updated"),
					testutils.CheckIPSetExists(IPSetSecond),
				),
			},
			// Update IPSet with multiple CIDRs
			{
				Config: testutils.HclIPSetResource(IPSetNameSecond, "ipset-comment-updated", IPSetCIDR),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", IPSetNameSecond),
					resource.TestCheckResourceAttr(tfNode, "comment", "ipset-comment-updated"),
					testutils.CheckIPSetExists(IPSetCIDR),
				),
			},
			// Update IPSet's include referenced alias
			{
				Config: testutils.HclIPSetWithAliasResource(IPSetNameSecond, "ipset-comment-updated", IPSetCIDRWithAlias),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", IPSetNameSecond),
					resource.TestCheckResourceAttr(tfNode, "comment", "ipset-comment-updated"),
					testutils.CheckIPSetExists(IPSetCIDRWithAlias),
				),
			},
		},
	})
}

// CheckIPSetDestroyed verifies that all ip sets referenced in the state
// are destroyed. This will be invoked *after* terraform destroys
// the resource but *before* the state is wiped clean
func CheckIPSetDestroyed(s *terraform.State) error {
	config := testutils.GetProvider().Meta().(proxmoxtf.ProviderConfiguration)
	conn, err := config.GetVEClient()

	if err != nil {
		return err
	}

	// loop through the resource state
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "proxmox_virtual_environment_cluster_ipset" {
			continue
		}

		response, err := conn.GetListIPSetContent(rs.Primary.ID)

		if err == nil {

			if len(response) != 0 {
				fmt.Errorf("IPSet with name `%s` should not exist", rs.Primary.ID)
			}

			return nil
		}

		if !strings.Contains(err.Error(), "no such IPSet") {
			return err
		}
	}

	return nil
}
