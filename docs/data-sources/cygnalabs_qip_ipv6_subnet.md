# [QIP IPv6 Subnet]

Use the `cygnalabs_qip_ipv6_subnet` data source to retrieve information for an IPv6 subnet managed by VitalQIP:

* `org_name` - `string`: **required**, Organization Name.
* `subnet_address` - `string`: **required**, IPv6 address of subnet.
* `subnet_prefix_length` - `int`: **required**, Prefix length of subnet.
* `subnet_name` - `string`: **optional**, Name of subnet.
* `pool_name` - `string`: **optional**, Name of pool.
* `block_name` - `string`: **optional**, Name of block.


### Example of a IPv6 Subnet

This example defines a data source of type `cygnalabs_qip_ipv6_subnet` with the name `ipv6_subnet_ds`, as configured in a Terraform file. By using this data source, you can reference and retrieve information about the specified IPv6 subnet.

```hcl
data "cygnalabs_qip_ipv6_subnet" "ipv6_subnet_ds" {
  org_name= "Terraform"
  subnet_address="2000::"
  subnet_prefix_length=56
}

output "ipv6_subnet_output" {
  value = data.cygnalabs_qip_ipv6_subnet.ipv6_subnet_ds 
}

```