package testutils

import (
	"bytes"
	"fmt"
	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/danitso/terraform-provider-proxmox/proxmoxtf"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"reflect"
	"text/template"
)

// CheckIPSetExists Given the name of an IPSet, this will return a function that will check
// whether or not an IPSet
// - (1) exists in the state
// - (2) exist in Proxmox VE
// - (3) has the correct name
func CheckIPSetExists(IPSetContent proxmox.VirtualEnvironmentClusterIPSetContent) resource.TestCheckFunc {
	return func (s *terraform.State) error {
		res, ok := s.RootModule().Resources["proxmox_virtual_environment_cluster_ipset.ipset"]
		if !ok {
			return fmt.Errorf("Did not find IPSet in the TF state")
		}

		clients := GetProvider().Meta().(proxmoxtf.ProviderConfiguration)
		id := res.Primary.ID
		ReadIPSetContent, err := readIPSet(clients, id)

		if err != nil {
			return fmt.Errorf("IPSet with Name=%s cannot be found. Error %v", id, err)
		}

		if len(IPSetContent) != len(ReadIPSetContent) {
			return fmt.Errorf("IPSet with Name=%s has content length %v, expected %v", id, len(ReadIPSetContent), len(IPSetContent))
		}

		if len(IPSetContent) == 0 && len(ReadIPSetContent) > 0 {
			return fmt.Errorf("IPSet with Name=%s should not have any elements in it. Found %v", id, len(ReadIPSetContent))
		}

		if len(IPSetContent) > 0 {

			for i, v := range IPSetContent {
				if v.CIDR != ReadIPSetContent[i].CIDR {
					return fmt.Errorf("IPSet with Name=%s contains not expected IP/CIDR=%s, expected %s", id, ReadIPSetContent[i].CIDR, v.CIDR)
				}

				if v.Comment != ReadIPSetContent[i].Comment {
					return fmt.Errorf("IPSet with Name=%s contains not expected Comment value=%s, exepected %s", id, ReadIPSetContent[i].Comment, v.Comment)
				}

				if !reflect.DeepEqual(&v.NoMatch, &ReadIPSetContent[i].NoMatch) {
					return fmt.Errorf("IPSet with Name=%s contains not expected NoMatch=%p, expected %p", id, ReadIPSetContent[i].NoMatch, v.NoMatch)
				}
			}

		}

		return nil
	}
}

// readIPSet is a helper function that reads an IPSet based on a given name
func readIPSet(clients proxmoxtf.ProviderConfiguration, identifier string) ([]*proxmox.VirtualEnvironmentClusterIPSetGetResponseData, error) {
	conn, err := clients.GetVEClient()

	if err != nil {
		return nil, err
	}

	response, err := conn.GetListIPSetContent(identifier)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// HclIPSetResource HCL describing of a PVE IPSet resource
func HclIPSetResource(name string, comment string, cidr proxmox.VirtualEnvironmentClusterIPSetContent) string  {
	var b bytes.Buffer

	if name == "" {
		panic("Parameter: `name` cannot be empty")
	}

	testRunData := map[string]interface{}{
		"Name" : name,
		"Comment": comment,
		"CIDR" : cidr,
	}

	tmpl, err := template.New("").Parse(`
resource "proxmox_virtual_environment_cluster_ipset" "ipset" {
	name    = "{{.Name}}"
	comment = "{{.Comment}}"
	{{range .CIDR}}
	ipset {
		cidr = "{{.CIDR}}"
		comment = "{{.Comment}}"
		{{if .NoMatch}}
		nomatch = "{{.NoMatch}}"
		{{end}}
	}
	{{end}}
}
`)
	t := template.Must(tmpl, err)
	if err := t.Execute(&b, testRunData); err != nil {
		_ = fmt.Errorf("Unable to parse template: %p", err)
		return ""
	}

	return b.String()
}

// HclIPSetWithAliasResource HCL describing of a PVE IPSet resource with alias
func HclIPSetWithAliasResource(name string, comment string, cidr proxmox.VirtualEnvironmentClusterIPSetContent) string  {
	alias := HclAliasResource("test-alias", "192.168.0.0/23", "alias-comment")
	IPSet := HclIPSetResource(name, comment, cidr)

	return fmt.Sprintf("%s\n%s", alias, IPSet)

}
