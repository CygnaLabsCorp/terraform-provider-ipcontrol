# Terraform Resource: IPC DNS RR

## Overview

The `ipcontrol_dns_rr` data source retrieves information about a DNS resource record managed by IPControl.

## Parameters

### Required Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `domain` | `string` | Domain name where resource records are to be added. |
| `owner` | `string` | The owner field of the resource record. |
| `resource_record_type` | `string` | The type of resource record being imported. To import a type not predefined in IPControl the type must be prepended with “Other”. |


### Optional Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `data` | `string` | The data portion of the resource record. The format is dependent on the type specified above. |

### Computed Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `domain_type` | `string` | The name of the domain type to which the domain belongs. |
| `resource_record_class` | `string` | The Class of the resource record. Defaults to “IN”. |
| `ttl` | `string` | Time to live. |
| `comment` | `string` | Comment text associated with the resource record. |
| `device_rec_flag` | `string` | When true, this indicates that the resource record is bound to a device. When false, this indicates that the resource record is associated with the domain only, and not a specific device. |



## Example Usage

#### Domain Data Source
```hcl
data "ipcontrol_dns_rr" "my_rr" {
  owner                = "caa"
  domain               = "com."
  resource_record_type = "A"
}

output "rr" {
  value = data.ipcontrol_dns_rr.my_rr
}

```

#### Domain Data Source with data
```hcl
data "ipcontrol_dns_rr" "my_rr" {
  owner                = "caa"
  domain               = "com."
  resource_record_type = "AAAA"
  data                 = "2001:db8:85a3::3000:82"
}

output "rr" {
  value = data.ipcontrol_dns_rr.my_rr
}

```