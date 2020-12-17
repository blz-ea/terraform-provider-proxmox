package testutils

import (
	"bytes"
	"fmt"
	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/danitso/terraform-provider-proxmox/proxmoxtf"
	"github.com/hashicorp/terraform/terraform"
	"html/template"
	"strconv"
)

type VMStateAttributes struct {
	NodeName 	string
	VMID 		int
}

// readVM is helper function that reads VM
func ReadVM(clients proxmoxtf.ProviderConfiguration, nodeName string, vmID int) (*proxmox.VirtualEnvironmentVMGetResponseData, error) {
	conn, err := clients.GetVEClient()

	if err != nil {
		return nil, err
	}

	res, err := conn.GetVM(nodeName, vmID)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetVMFromState helper function that retrieves VM Attributes from TF state
func GetVMFromState(s *terraform.State, resourceName string) (VMStateAttributes, error) {

	vm, ok := s.RootModule().Resources[fmt.Sprintf("proxmox_virtual_environment_vm.%s", resourceName)]

	if !ok {
		return VMStateAttributes{}, fmt.Errorf("Did not find a VM with name %s in TF state", resourceName)
	}

	nodeName := vm.Primary.Attributes["node_name"]
	vmID, err := strconv.Atoi(vm.Primary.Attributes["vm_id"])

	if err != nil {
		return VMStateAttributes{}, fmt.Errorf("Unable to convert `vm_id` attribute to integer")
	}

	vmAttributes := VMStateAttributes{
		NodeName: nodeName,
		VMID: vmID,
	}

	return vmAttributes, nil
}

// HclVMResource HCL describing of a PVE VM resource
func HclVMResource(vm map[string]string) string {
	var b bytes.Buffer

	tmpl, err := template.New("").Parse(`
{{ if .PoolID }}
resource "proxmox_virtual_environment_pool" "test_pool" {
  comment = "Managed by Terraform"
  pool_id = "{{ .PoolID }}"
}
{{ end }}

data "proxmox_virtual_environment_nodes" "pve_nodes" {}

resource "proxmox_virtual_environment_vm" "vm" {
  {{ if .NodeName }}
  node_name   = "{{.NodeName}}"
  {{ else }}
  node_name   = data.proxmox_virtual_environment_nodes.pve_nodes.names[0]
  {{ end }}

  {{ if .Name }}
  name = "{{.Name}}"
  {{ end }}

  {{ if .PoolID }}
  pool_id     = proxmox_virtual_environment_pool.test_pool.pool_id
  {{ end }}
}
`)

	t := template.Must(tmpl, err)
	if err := t.Execute(&b, vm); err != nil {
		_ = fmt.Errorf("Unable to parse template %p", err)
		return ""
	}

	return b.String()
}