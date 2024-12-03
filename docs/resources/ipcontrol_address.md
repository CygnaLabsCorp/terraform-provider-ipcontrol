# Terraform Resource: IPC Address

## Overview

The resource manages device configurations with comprehensive network and interface settings.

## Parameters

### Required Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `hostname` | `string` | The host name of the Device. Requires creating a new resource if changed. |
| `interfaces` | `list` | Network interfaces for the device. Each interface requires specific details. |

### Optional Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `domain_name` | `string` | The Domain name of the Device. |
| `device_type` | `string` | The type of the device. |
| `domain_type` | `string` | The type of the domain. |
| `duid` | `string` | The DHCP Unique Identifier for IPv6. |
| `options` | `list(string)` | Additional configuration options. |

### Computed Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `string` | Unique identifier for the Device. |
| `container` | `string` | The name of the container that contains the device. |
| `ip_address` | `string` | The IP Address of the Device. |
| `address_type` | `string` | The address type of this interface IP address. |
| `resource_record_flag` | `string` | Resource record status. |
| `aliases` | `list(string)` | Device aliases. |

### Interface Configuration

Each interface in the `interfaces` list supports:

| Parameter | Type | Description |
|-----------|------|-------------|
| `address_type` | `list(string)` | Type of IP address (Required) |
| `ip_address` | `list(string)` | IP addresses for the interface (Required) |
| `container` | `list(string)` | Container information (Optional, Computed) |
| `name` | `string` | Interface name (Optional, Computed) |
| `id` | `integer` | Interface identifier (Computed) |

## Supported Address Type
#### For IPv4
- Static
- Dynamic DHCP
- Automatic DHCP
- Manual DHCP
- Reserved
#### For IPv6
- Dynamic NA DHCPv6
- Automatic NA DHCPv6
- Manual NA DHCPv6
- Dynamic TA DHCPv6
- Automatic TA DHCPv6

> Note: Note that if Dynamic, Automaticor Manual DHCP is specified, there must be a DHCP server defined in the subnet policies for this IP Address

### Options Configuration

Each item in the `options` list supports:

| Value | Description |
|-----------|-------------|
| `ignoreDupWarning` | TWhen this option is specified, if the administrator policy of the user indicates "warn" for the "Allow Duplicate Hostname Checking" option, the warning will be ignored and the device added with the duplicate hostname. |
| `resourceRecordFlag` | When this option is specified, resource records are added for this device. |
| `splitPool` | When this option is specified, if the IPv4 address specified in column A falls within an existing address pool, the device is created and the pool is split around the given IP address. This function is not supported for IPv6. |
| `contiguous` | When this option is specified, the number of addresses specified in rangeSize will be allocated as a contiguous IP address range. Otherwise, the rangeSize number of devices will be added, but may not be in a contiguous range |
| `rangeSize` | Specify the number of addresses to be added as rangeSize=number, for example, rangeSize=10. This defaults to 1. |
| `instantDNS` | Create an immediate "DNS Configuration - Chained Resource Records Only" (via DDNS) task for auto-created device resource records. Applicable only if the resourceRecordFlag is specified. |

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