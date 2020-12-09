/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

/**
* Reference: https://pve.proxmox.com/pve-docs/api-viewer/#/cluster/firewall/ipset
 */

package proxmox

// GET /api2/json/cluster/firewall/ipset
// VirtualEnvironmentClusterIPSetListResponseBody contains the data from an IPSet get response.
type VirtualEnvironmentClusterIPSetListResponseBody struct {
	Data []*VirtualEnvironmentClusterIPSetCreateRequestBody `json:"data,omitempty"`
}

// POST /api2/json/cluster/firewall/ipset
// VirtualEnvironmentClusterIPSetCreateRequestBody contains the data for an IPSet create request
type VirtualEnvironmentClusterIPSetCreateRequestBody struct {
	Comment string `json:"comment,omitempty" url:"comment,omitempty"`
	Name    string  `json:"name" url:"name"`
}

// GET /api2/json/cluster/firewall/ipset/{name}
// VirtualEnvironmentClusterIPSetGetResponseBody contains the body from an IPSet get response.
type VirtualEnvironmentClusterIPSetGetResponseBody struct {
	Data []*VirtualEnvironmentClusterIPSetGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentClusterIPSetGetResponseData contains the data from an IPSet get response.
type VirtualEnvironmentClusterIPSetGetResponseData struct {
	CIDR		string `json:"cidr" url:"cidr"`
	NoMatch		*CustomBool  `json:"nomatch,omitempty" url:"nomatch,omitempty,int"`
	Comment		string `json:"comment,omitempty" url:"comment,omitempty"`
}

// POST /api2/json/cluster/firewall/ipset
// VirtualEnvironmentClusterIPSetUpdateRequestBody contains the data for an IPSet update request.
type VirtualEnvironmentClusterIPSetUpdateRequestBody struct {
	ReName  string `json:"rename,omitempty" url:"rename,omitempty"`
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	Name    string  `json:"name" url:"name"`
}

// GET /api2/json/cluster/firewall/ipset
// VirtualEnvironmentClusterIPSetGetResponseData contains list of IPSets from
type VirtualEnvironmentClusterIPSetListResponseData struct {
	Comment		*string `json:"comment,omitempty" url:"comment,omitempty"`
	Name		string  `json:"name" url:"name"`
}

// VirtualEnvironmentClusterIPSetContent is an array of VirtualEnvironmentClusterIPSetGetResponseData.
type VirtualEnvironmentClusterIPSetContent []VirtualEnvironmentClusterIPSetGetResponseData
