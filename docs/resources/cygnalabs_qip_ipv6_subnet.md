# Resource: [QIP IPv6 Subnet]

###  Descriptions
The `cygnalabs_qip_ipv6_subnet` resource is designed to manage and associate an IPv6 subnet with an organizational structure within VitalQIP.

### Parameters
The following list describes the parameters you can define in the resource block of the record:

* `org_name` - `string`: **required**, Organization Name.
* `subnet_address` - `string`: **required**, IPv6 address of subnet.
* `subnet_prefix_length` - `int`: **required**, Prefix length of subnet.
* `block_address` - `string`: **required**, IPv6 address of block.
* `block_prefix_length` - `int`: **required**, Prefix length of block.
* `create_reverse_zone` - `boolean`: **required**, Create Reverse Zone.
* `subnet_name` - `string`: **optional**, Name of subnet.
* `pool_name` - `string`: **optional**, Name of pool.
* `block_name` - `string`: **optional**, Name of block.

### ⚠️ Force replacement fields
The following fields after changes will require deleting and recreating the resource:
* `org_name` - Can't change after created in VitalQIP.
* `subnet_address` - Can't change after created in VitalQIP.
* `subnet_prefix_length` - Can't change after created in VitalQIP.
* `block_address` - Can't change after created in VitalQIP.
* `block_prefix_length` - Can't change after created in VitalQIP.

> **WARNING**: Changing the above fields will result in the current resource being deleted and a new one created. Make sure you back up your data and understand the impact before making changes.

### Ignoring changes to fields

These fields are only used when creating an IPv6 subnet. The API for updating IPv6 subnets does not support changing these fields, so they will be ignored.
* `pool_name`
* `block_name`
* `create_reverse_zone`

## How to use
First define `resource` in the .tf file.<br>
`IPv6 Subnet` example
```hcl
resource "cygnalabs_qip_ipv6_subnet" "ipv6_subnet" {
  // required parameters
  org_name= "Terraform"
  subnet_address="2000:0:0:10::"
  subnet_name="subnet_name"
  subnet_prefix_length = 60
  block_prefix_length = 48
  block_address="2000::"
  create_reverse_zone=true
  
  // optional parameters
  pool_name="pool_name"
  block_name="block_name"

  // ignoring changes to fields
  lifecycle {
    ignore_changes = [pool_name, block_name, create_reverse_zone]
  }
}

output "export_ipv6_subnet" {
  value = resource.cygnalabs_qip_ipv6_subnet.ipv6_subnet
}

```

Then run
```bash
terraform apply
```