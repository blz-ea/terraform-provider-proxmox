/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmoxtf

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	mkDataSourceVirtualEnvironmentPoolsPoolIDs = "pool_ids"
)

func dataSourceVirtualEnvironmentPools() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			mkDataSourceVirtualEnvironmentPoolsPoolIDs: {
				Type:        schema.TypeList,
				Description: "The pool ids",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Read: dataSourceVirtualEnvironmentPoolsRead,
	}
}

func dataSourceVirtualEnvironmentPoolsRead(d *schema.ResourceData, m interface{}) error {
	config := m.(ProviderConfiguration)
	veClient, err := config.GetVEClient()

	if err != nil {
		return err
	}

	list, err := veClient.ListPools()

	if err != nil {
		return err
	}

	poolIDs := make([]interface{}, len(list))

	for i, v := range list {
		poolIDs[i] = v.ID
	}

	d.SetId("pools")

	d.Set(mkDataSourceVirtualEnvironmentPoolsPoolIDs, poolIDs)

	return nil
}
