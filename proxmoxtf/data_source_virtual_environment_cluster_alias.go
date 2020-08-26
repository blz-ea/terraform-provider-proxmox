/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"github.com/danitso/terraform-provider-proxmox/proxmox"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

const (
	dvResourceVirtualEnvironmentAliasComment 	= ""

	mkResourceVirtualEnvironmentAliasComment    = "comment"
	mkResourceVirtualEnvironmentAliasName 		= "name"
	mkResourceVirtualEnvironmentAliasCIDR 		= "cidr"
)

func resourceVirtualEnvironmentAlias() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkResourceVirtualEnvironmentAliasComment: {
				Type: schema.TypeString,
				Description: "Alias comment",
				Optional: true,
				Default: dvResourceVirtualEnvironmentAliasComment,
			},
			mkResourceVirtualEnvironmentAliasCIDR: {
				Type: schema.TypeString,
				Description: "IP/CIDR block",
				Required: true,
				ForceNew: false,
			},
			mkResourceVirtualEnvironmentAliasName: {
				Type: schema.TypeString,
				Description: "Alias name",
				Required: true,
				ForceNew: false,
			},
		},
		Create: resourceVirtualEnvironmentAliasCreate,
		Read: resourceVirtualEnvironmentAliasRead,
		Update: resourceVirtualEnvironmentAliasUpdate,
		Delete: resourceVirtualEnvironmentAliasDelete,
	}
}

func resourceVirtualEnvironmentAliasCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	comment := d.Get(mkResourceVirtualEnvironmentAliasComment).(string)
	name := d.Get(mkResourceVirtualEnvironmentAliasName).(string)
	cidr := d.Get(mkResourceVirtualEnvironmentAliasCIDR).(string)

	body := &proxmox.VirtualEnvironmentAliasCreateRequestBody{
		Comment: &comment,
		Name: name,
		CIDR: cidr,
	}

	err = veClient.CreateAlias(body)

	if err != nil {
		return err
	}

	d.SetId(name)

	return resourceVirtualEnvironmentAliasRead(d, m)
}

func resourceVirtualEnvironmentAliasRead(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return nil
	}

	name := d.Id()
	alias, err := veClient.GetAlias(name)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")
			return nil
		}

		return err
	}

	aliasMap := map[string]interface{}{
		mkResourceVirtualEnvironmentAliasComment: alias.Comment,
		mkResourceVirtualEnvironmentAliasName: alias.Name,
		mkResourceVirtualEnvironmentAliasCIDR: alias.CIDR,
	}

	for key, val := range aliasMap {
		err = d.Set(key, val)

		if err != nil {
			return err
		}
	}

	return nil
}

func resourceVirtualEnvironmentAliasUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	comment := d.Get(mkResourceVirtualEnvironmentAliasComment).(string)
	cidr := d.Get(mkResourceVirtualEnvironmentAliasCIDR).(string)
	newName := d.Get(mkResourceVirtualEnvironmentAliasName).(string)
	previousName := d.Id()

	body := &proxmox.VirtualEnvironmentAliasUpdateRequestBody{
		ReName: newName,
		CIDR: cidr,
		Comment: &comment,
	}

	err = veClient.UpdateAlias(previousName, body)

	if err != nil {
		return err
	}

	d.SetId(newName)

	return resourceVirtualEnvironmentAliasRead(d, m)
}


func resourceVirtualEnvironmentAliasDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(providerConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return nil
	}

	name := d.Id()
	err = veClient.DeleteAlias(name)

	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") {
			d.SetId("")
			return nil
		}

		return err
	}

	d.SetId("")

	return nil
}














