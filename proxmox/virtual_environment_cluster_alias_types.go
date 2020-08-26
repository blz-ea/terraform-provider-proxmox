/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package proxmox

// VirtualEnvironmentAliasCreateRequestBody contains the data for an alias create request.
type VirtualEnvironmentAliasCreateRequestBody struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	Name    string  `json:"name" url:"name"`
	CIDR    string  `json:"cidr" url:"cidr"`
}

// VirtualEnvironmentAliasGetResponseBody contains the body from an alias get response.
type VirtualEnvironmentAliasGetResponseBody struct {
	Data *VirtualEnvironmentAliasGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentAliasGetResponseData contains the data from an alias get response.
type VirtualEnvironmentAliasGetResponseData struct {
	Comment		*string `json:"comment,omitempty" url:"comment,omitempty"`
	Name		string  `json:"name" url:"name"`
	CIDR		string  `json:"cidr" url:"cidr"`
	Digest  	*string  `json:"digest" url:"digest"`
	IPVersion	int		`json:"ipversion" url:"ipversion"`
}

// VirtualEnvironmentAliasGetResponseData contains the data from an alias get response.
type VirtualEnvironmentAliasListResponseBody struct {
	Data []*VirtualEnvironmentAliasGetResponseData `json:"data,omitempty"`
}

// VirtualEnvironmentAliasUpdateRequestBody contains the data for an alias update request.
type VirtualEnvironmentAliasUpdateRequestBody struct {
	Comment *string `json:"comment,omitempty" url:"comment,omitempty"`
	ReName	string  `json:"rename" url:"rename"`
	CIDR	string  `json:"cidr" url:"cidr"`
}

