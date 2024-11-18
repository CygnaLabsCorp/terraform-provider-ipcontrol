# Resource: [QIP IPv4 Subnet]

###  Descriptions
The `cygnalabs_qip_ipv4_subnet` resource is designed to manage and associate an IPv4 subnet with an organizational structure within VitalQIP.

### Parameters
The following list describes the parameters you can define in the resource block of the record:

* `org_name` - `string`: **required**, Organization Name.
* `subnet_address` - `string`: **required**, IPv4 address of subnet.
* `network_address` - `string`: **required**, IPv4 address of network.
* `subnet_mask` - `string`: **required**, Subnet mask of subnet.
* `warning_percent` - `int`: **required**, Percentage of managed addresses. Warning if this percentage is reached. The value range is 0 - 100.
* `subnet_name` - `string`: **optional**, Name of subnet.
* `warning_type` - `int`: **optional**, Type of warning when the defined percentage of addresses reached.

### ⚠️ Force replacement fields
The following fields after changes will require deleting and recreating the resource:
* `org_name` - Can't change after created in VitalQIP.
* `subnet_address` - Can't change after created in VitalQIP.
* `network_address` - Can't change after created in VitalQIP.
* `subnet_mask` - Can't change after created in VitalQIP.

> **WARNING**: Changing the above fields will result in the current resource being deleted and a new one created. Make sure you back up your data and understand the impact before making changes.

## How to use
First define `resource` in the .tf file.<br>
`IPv4 Subnet` example
```hcl
resource "cygnalabs_qip_ipv4_subnet" "ipv4_subnet" {
  // required parameters
  org_name= "Terraform"
  subnet_address="10.0.0.0"
  subnet_mask = "255.255.255.0"
  network_address="10.0.0.0"
  warning_percent=90
  
  // optional parameters
  subnet_name="subnet_name"
  warning_type=0
}

output "ipv4_subnet_output" {
  value = resource.cygnalabs_qip_ipv4_subnet.ipv4_subnet
}
```

Then run
```bash
terraform apply
```