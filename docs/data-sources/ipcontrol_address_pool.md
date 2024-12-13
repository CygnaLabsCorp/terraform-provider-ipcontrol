# Terraform Data Source: IPC Address Pool

## Overview


The `ipcontrol_address_pool` data source retrieves information about an address pool managed by IPControl.

## Parameters

### Required Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `start_address` | `string` | The start of the address pool. |

### Optional Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `container` | `string` | The name of the container that holds the block in which the pool is defined. This is required only if there is overlapping address space in use and the start address is in overlapping space. |

### Computed Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `type` | `string` | The type of address pool. |
| `dhcp_policy_set` | `string` | The name of a Policy Set used with this pool. |
| `dhcp_option_set` | `string` | The name of a Option Set used with this pool. |
| `overlap_interface_ip` | `boolean` | Flag to allow a DHCPv6 pool to overlap an interface address. |
| `name` | `string` | The name of the address pool. |
| `last_admin` | `string` | The login id of the last administrator to update the pool. |
| `create_date` | `string` | The date the pool was created. |
| `end_address` | `string` | The end of the address pool.|
| `prefix_length` | `number` | The size parameter represents the subnet mask or prefix length of the address block in CIDR notation. Required for IPv6, the size value is larger due to the increased address space. IPv6 prefix lengths commonly range between /48 to /128, with /64 often used as the standard size for a single subnet.|
| `primary_net_service` | `string` |The name of the DHCP server that will serve addresses from this pool.|



#### Address Pool with IPv4
```hcl
data "ipcontrol_address_pool" "my_pool_v4" {
  start_address = "10.0.0.0"

  # optional parameter
  container = "InControl/caa"
}

output "export_01" {
  value = data.ipcontrol_address_pool.my_pool_v4
}
```
#### Address Pool with IPv6
```hcl
data "ipcontrol_address_pool" "my_pool_v6" {
  start_address = "2001:db8:85a3::3000:0"
}

output "export_02" {
  value = data.ipcontrol_address_pool.my_pool_v6
}

```