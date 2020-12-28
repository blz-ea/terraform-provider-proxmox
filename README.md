# Terraform Provider for Proxmox

- Website: [terraform.io](https://registry.terraform.io/providers/blz-ea/proxmox/)

A Terraform Provider which adds support for Proxmox solutions.
 
<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="300px">

## Requirements
- [Terraform](https://www.terraform.io/downloads.html) 0.13+
- [Go](https://golang.org/doc/install) 1.15 (to build the provider plugin)

## Building the Provider
Clone repository to: `$GOPATH/src/github.com/blz-ea/terraform-provider-proxmox`

```sh
$ mkdir -p $GOPATH/src/github.com/blz-ea; cd $GOPATH/src/github.com/blz-ea
$ git clone git@github.com:blz-ea/terraform-provider-proxmox
```

Since Terraform 0.13, all providers (including local ones) require a version and need to meet stricter load conditions. To simplify this there is a make target available which will handle building and putting the executable in the correct place.
```sh
make build-and-install-dev-version
```

Once this completes, you will need to update your Terraform provider configuration to reference version 99.0.0 (the version of the development build) in order to use it. Example configuration:

```hcl
terraform {
  required_version = ">= 0.13"
  required_providers {
    proxmox = {
      source = "blz-ea/proxmox"
      version = "99.0.0"
    }
  }
}
```

Be sure to remove any additional configuration from your `~/.terraformrc` and `.terraform` directories to prevent load issues. You will also need to update all `version` references to use 99.0.0 in order to fully initialise the development provider. If you see other versions (like the example below) you need to track them down and remove them.

```
Finding blz-ea/proxmox versions matching "99.0.0, ~> 0.*"...
```

## Developing the Provider
If you wish to work on the provider, you'll first need [Go](http://www.golang.org)
installed on your machine (version 1.15+ is *required*). You'll also need to
correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well
as adding `$GOPATH/bin` to your `$PATH`.

See above for which option suits your workflow for building the provider.

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests

**Note:** Acceptance tests provision real resources

Acceptance tests require running Proxmox Virtual Environment.
To run ant acceptance tests you need to set `PROXMOX_VE_ENDPOINT`, `PROXMOX_VE_USERNAME`, 
`PROXMOX_VE_PASSWORD`, `PROXMOX_VE_INSECURE` environmental variables.

```sh
export PROXMOX_VE_ENDPOINT='<protocol>://<host>:<port>'
export PROXMOX_VE_USERNAME='<username>'
export PROXMOX_VE_PASSWORD='<password>'
export PROXMOX_VE_INSECURE='1'

$ make testacc
```

or create `.env` file in the root of the folder with the following content

```env
PROXMOX_VE_ENDPOINT='<protocol>://<host>:<port>'
PROXMOX_VE_USERNAME='<username>'
PROXMOX_VE_PASSWORD='<password>'
PROXMOX_VE_INSECURE='1'
```

```sh
$ make testacc
```

If you wish to contribute to the provider, please see [CONTRIBUTING.md](CONTRIBUTING.md).

## Known issues

### Disk images cannot be imported by non-PAM accounts
Due to limitations in the Proxmox VE API, certain actions need to be performed using SSH. This requires the use of a PAM account (standard Linux account).

### Disk images from VMware cannot be uploaded or imported
Proxmox VE is not currently supporting VMware disk images directly. However, you can still use them as disk images by using this workaround:

```hcl
resource "proxmox_virtual_environment_file" "vmdk_disk_image" {
  content_type = "iso"
  datastore_id = "datastore-id"
  node_name    = "node-name"

  source_file {
    # We must override the file extension to bypass the validation code in the Proxmox VE API.
    file_name = "vmdk-file-name.img"
    path      = "path-to-vmdk-file"
  }
}

resource "proxmox_virtual_environment_vm" "example" {
  ...

  disk {
    datastore_id = "datastore-id"
    # We must tell the provider that the file format is vmdk instead of qcow2.
    file_format  = "vmdk"
    file_id      = "${proxmox_virtual_environment_file.vmdk_disk_image.id}"
  }

  ...
}
```

### Snippets cannot be uploaded by non-PAM accounts
Due to limitations in the Proxmox VE API, certain files need to be uploaded using SFTP. This requires the use of a PAM account (standard Linux account).
