# Terraform Resource: IPC Address

## Overview

The `ipcontrol_address` data source retrieves information about a device managed by IPControl.

## Parameters

### Required Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `ip_address` | `string` | The IP Address of the Device. |


### Computed Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `string` | Unique identifier for the Device. |
| `container` | `string` | The name of the container that contains the device. |
| `ip_address` | `string` | The IP Address of the Device. |
| `address_type` | `string` | The address type of this interface IP address. |
| `resource_record_flag` | `string` | Resource record status. |
| `aliases` | `list(string)` | Device aliases. |
| `domain_name` | `string` | The Domain name of the Device. |
| `hostname` | `string` | The host name of the Device. |
| `device_type` | `string` | The type of the device. |
| `domain_type` | `string` | The type of the domain. |
| `duid` | `string` | The DHCP Unique Identifier for IPv6. |

### Interface Configuration

Each interface in the `interfaces` list supports:

| Parameter | Type | Description |
|-----------|------|-------------|
| `address_type` | `list(string)` | Type of IP address (Computed) |
| `ip_address` | `list(string)` | IP addresses for the interface (Computed) |
| `container` | `list(string)` | Container information (Computed) |
| `name` | `string` | Interface name (Computed) |
| `id` | `integer` | Interface identifier (Computed) |


## Example Usage

#### Device with IPv4
```hcl
resource "ipcontrol_address" "pc_v4" {
  options = [
    "ignoreDupWarning",
    "resourceRecordFlag"
  ]
  
  device_type = "PC"
  domain_name = "com."
  hostname    = "tfhost"

  interfaces {
    name         = "tfname"
    address_type = ["Static"]
    ip_address   = ["135.0.0.146"]
  }
}
```
#### Device with IPv6
```hcl
resource "ipcontrol_address" "example" {
  options = [
    "ignoreDupWarning",
    "resourceRecordFlag"
  ]
  
  device_type = "PC"
  domain_name = "com."
  hostname    = "tfhostv6"

  interfaces {
    name         = "tfname"
    address_type = ["Static"]
    ip_address   = ["2001:db8:85a3::3000:0"]
  }
}
```