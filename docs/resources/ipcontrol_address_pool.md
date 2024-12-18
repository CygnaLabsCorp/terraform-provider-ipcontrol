# Terraform Resource: IPC Address Pool

## Overview

The `ipcontrol_address_pool` resource manages address pool in IPControl.

## Parameters

### Required Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `start_address` | `string` | The start of the address pool. |
| `type` | `string` | The type of address pool. |

### Required Parameters with conditional

| Parameter | Type | Description |
|-----------|------|-------------|
| `end_address` | `string` | The end of the address pool. **Required** for IPv4, IPv6 will be computed |
| `prefix_length` | `number` | The size parameter represents the subnet mask or prefix length of the address block in CIDR notation. Required for IPv6, the size value is larger due to the increased address space. IPv6 prefix lengths commonly range between /48 to /128, with /64 often used as the standard size for a single subnet. **Required** for IPv6, IPv4 ignore |
| `primary_net_service` | `string` |The name of the DHCP server that will serve addresses from this pool. Note: This is **Required** when a DHCP server is not defined for the subnet, and this address pool is one of the DHCP address types:  “Dynamic DHCP”, “Automatic DHCP”, “Dynamic NA DHCPv6”,  Automatic NA DHCPv6”, ”Dynamic TA DHCPv6”, ”Automatic TA DHCPv6” |

### Optional Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `container` | `string` | The name of the container that holds the block in which the pool is defined. This is required only if there is overlapping address space in use and the start address is in overlapping space. |
| `dhcp_policy_set` | `string` | The name of a Policy Set used with this pool. |
| `dhcp_option_set` | `string` | The name of a Option Set used with this pool. |
| `overlap_interface_ip` | `boolean` | Flag to allow a DHCPv6 pool to overlap an interface address. |
| `name` | `string` | The name of the address pool. |

### Computed Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `last_admin` | `string` | The login id of the last administrator to update the pool. |
| `create_date` | `string` | The date the pool was created. |



## Supported Type
#### For IPv4
- Static
- Dynamic DHCP
- Automatic DHCP
- Reserved
#### For IPv6
- Dynamic NA DHCPv6
- Dynamic TA DHCPv6
- Automatic NA DHCPv6
- Automatic TA DHCPv6
- Reserved

## Example Usage

#### Address Pool with IPv4
```hcl
resource "ipcontrol_address_pool" "my_pool" {
  start_address       = "135.0.0.161"
  end_address         = "135.0.0.164"
  name                = "my-addrp"
  type                = "Dynamic DHCP"
  primary_net_service = "dhcp"

  lifecycle {
    ignore_changes = [overlap_interface_ip, prefix_length]
  }
}
```
#### Address Pool with IPv6
```hcl
resource "ipcontrol_address_pool" "my_pool_v6" {
  start_address        = "2001:db8:85a3::3000:0"
  name                 = "my-addrp2"
  type                 = "Dynamic NA DHCPv6"
  prefix_length        = 128
  primary_net_service  = "dhcpv6"
  dhcp_option_set      = "Cisco DHCPv6 Option Set"
  dhcp_policy_set      = "Cisco DHCP 8.0 Client Class Template Policy Set"
  overlap_interface_ip = true
}

```