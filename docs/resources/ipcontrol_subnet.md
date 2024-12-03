# Terraform Resource: ipcontrol_subnet

## Overview

The `ipcontrol_subnet` resource is used to assign IP address blocks for devices within a specified container.

## Parameters

### Required Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `container` | `string` | The name of the container that will hold the block. |
| `address` | `string` | The IP address block to allocate. |
| `size` | `int` | The subnet mask or prefix length of the address block in CIDR notation. |

### Optional Parameters

| Parameter | Type | Description | Default |
|-----------|------|-------------|---------|
| `address_version` | `int` | IP Address version (4 for IPv4, 6 for IPv6) | 4 |
| `rawcontainer` | `boolean` | Pass container parameter through to API without prefixing | `false` |
| `type` | `string` | Block Type (defaults to "Any" if not specified) | "Any" |
| `dns_domain` | `string` | DNS domain name for the block | - |
| `name` | `string` | Name of the block | - |
| `block_status` | `string` | Current status of the block | - |
| `cloud_type` | `string` | Cloud Provider type | - |
| `cloud_object_id` | `string` | Object ID in the cloud environment | - |

## Address Version Details

### IPv4
- Size range: 0-32
- Example: 24 represents 255.255.255.0 subnet mask

### IPv6
- Size range: 48-128
- Typical standard size: /64

## Supported Cloud Providers

- AWS
- Azure
- Cisco ACI
- Cisco DNA Center
- CloudBolt
- OpenStack
- ServiceNow
- VMware

### ⚠️ Force Replacement Fields
The following fields after changes will require deleting and recreating the resource:
* `container` - Can't change after created in IPControl.
* `address` - Can't change after created in IPControl.
* `size` - Can't change after created in IPControl.

> **WARNING**: Changing the above fields will result in the current resource being deleted and a new one created. Make sure you back up your data and understand the impact before making changes.

## Example Usage

```hcl
resource "ipcontrol_subnet" "example" {
  container = "MyContainer"
  address   = "192.168.1.0"
  size      = 24
  type      = "Any"
}