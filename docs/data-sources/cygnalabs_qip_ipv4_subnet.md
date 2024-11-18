# [QIP IPv4 Subnet]

Use the `cygnalabs_qip_ipv4_subnet` data source to retrieve information for an IPv4 subnet managed by VitalQIP:

* `org_name` - `string`: **required**, Organization Name.
* `subnet_address` - `string`: **required**, IPv4 address of subnet.
* `network_address` - `string`: **optional**, IPv4 address of network.
* `subnet_mask` - `string`: **optional**, Subnet mask of subnet.
* `subnet_name` - `string`: **optional**, Name of subnet.
* `warning_percent` - `int`: **optional**, Percentage of managed addresses. Warning if this percentage is reached. The value range is 0 - 100.
* `warning_type` - `int`: **optional**, Type of warning when the defined percentage of addresses reached.


### Example of a IPv4 Subnet

This example defines a data source of type `cygnalabs_qip_ipv4_subnet` with the name `ipv4_subnet_ds`, as configured in a Terraform file. By using this data source, you can reference and retrieve information about the specified IPv4 subnet.

```hcl
data "cygnalabs_qip_ipv4_subnet" "ipv4_subnet_ds" {
  org_name= "Terraform"
  subnet_address="10.0.0.0"
}

output "ipv4_subnet_output" {
  value = data.cygnalabs_qip_ipv4_subnet.ipv4_subnet_ds 
}

```