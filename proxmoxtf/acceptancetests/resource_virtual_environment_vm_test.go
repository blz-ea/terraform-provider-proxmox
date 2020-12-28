package acceptancetests

import (
	"fmt"
	"github.com/danitso/terraform-provider-proxmox/proxmoxtf"
	"github.com/danitso/terraform-provider-proxmox/proxmoxtf/acceptancetests/testutils"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"
)

// Verifies that VM can be created and updated
func TestAccResourceVirtualEnvironmentVM_CreateAndUpdate(t *testing.T) {
	// Set environmental variable that will signal Terraform to
	// Stop the VM rather than wait for shutdown
	os.Setenv("TF_ACC_VM_FORCE_STOP", "true")

	tfNode := "proxmox_virtual_environment_vm.vm"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckVMDestroyed,
		Steps: []resource.TestStep{
			// Create empty VM
			{
				Config: testutils.HclVMResource(map[string]string{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", ""),
				),
			},
			// Update VM
			// - Add name
			{
				Config: testutils.HclVMResource(map[string]string{
					"Name": "AccTestVM",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", "AccTestVM"),
				),
			},
			// Update VM
			// - Remove name
			{
				Config: testutils.HclVMResource(map[string]string{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", ""),
				),
			},
			// Update VM
			// - Add name
			// - Add to Pool
			{
				Config: testutils.HclVMResource(map[string]string{
					"Name":   "AccTestVM",
					"PoolID": "test-pool",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", "AccTestVM"),
				),
			},
		},
	})
}

// CheckVMDestroyed verifies that vm referenced in the state
// was destroyed. This will be invoked *after* terraform destroys
// the resource but *before* the state is wiped clean
func CheckVMDestroyed(s *terraform.State) error {
	config := testutils.GetProvider().Meta().(proxmoxtf.ProviderConfiguration)
	vmState, err := testutils.GetVMFromState(s, "vm")

	if err != nil {
		return err
	}

	if vmState.NodeName == "" || vmState.VMID == 0 {
		return fmt.Errorf("Unable to find `proxmox_virtual_environment_vm` resource")
	}

	conn, err := config.GetVEClient()

	if err != nil {
		return err
	}

	_, err = conn.GetVM(vmState.NodeName, vmState.VMID)

	if err == nil {
		return fmt.Errorf("VM with id `%d` should not exist", vmState.VMID)
	}

	return nil
}
