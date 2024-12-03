# Terraform Data Source: ipcontrol_subnet

## Overview

The `ipcontrol_subnet` data source retrieves information about a block managed by IPControl.

## Parameters

### Required Parameters

| Parameter      | Type     | Description                                                                 |
|----------------|----------|-----------------------------------------------------------------------------|
| `container`    | `string` | The name of the container that will hold the block.                        |
| `address`      | `string` | The address block to allocate.                                             |
| `size`         | `int`    | The subnet mask or prefix length of the address block in CIDR notation.    |

### Optional Parameters

| Parameter          | Type      | Description                                                                                      | Default   |
|--------------------|-----------|--------------------------------------------------------------------------------------------------|-----------|
| `address_version`  | `int`     | IP Address version (4 for IPv4, 6 for IPv6).                                                     | 4         |
| `rawcontainer`     | `boolean` | Pass the container parameter through to the API without prefixing.                              | `false`   |
| `type`             | `string`  | Block Type (defaults to "Any" if not specified).                                                | "Any"    |
| `name`             | `string`  | The name of the block.                                                                           | -         |
| `block_status`     | `string`  | The current status of the block. Accepted values: `Deployed`, `FullyAssigned`, `Reserved`, `Aggregate`. | -         |
| `cloud_type`       | `string`  | Specify the type of Cloud Provider. Supported values: AWS, Azure, Cisco ACI, Cisco DNA Center, CloudBolt, OpenStack, ServiceNow, VMware. | -         |
| `cloud_object_id`  | `string`  | The ID of this object as it is known in the cloud environment.                                   | -         |

## Address Version Details

### IPv4
- Size range: 0-32
- Example: 24 represents 255.255.255.0 subnet mask

### IPv6
- Size range: 48-128
- Typical standard size: /64

## Example Usage

This example defines a data source of type `ipcontrol_subnet` named `my_ipc_ds`, which is configured in a Terraform file.
You can reference this data source to retrieve information about the block.

```hcl
data "ipcontrol_subnet" "my_ipc_ds" {
  container      = "InControl/caa"
  rawcontainer   = true
  address        = "10.0.0.0"
  address_version = 4
  size           = 25
}

// Accessing individual fields in results
output "my-ipc-ds" {
  value = data.ipcontrol_subnet.my_ipc_ds.address
}
```

