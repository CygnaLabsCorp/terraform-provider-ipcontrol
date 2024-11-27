# [IPC Subnet]

Use the `cygnalabs_ipc_subnet` data source to retrieve the following information for a block, which is managed by a IPControl:

* `container` - `string`: **required**, The name of the container that will hold the block.
* `address` - `string`: **required**, The address block to allocate.
* `size` - `int`: **required**, The size parameter represents the subnet mask or prefix length of the address block in CIDR notation. For IPv4, this is typically a value between 0 and 32 (e.g., 24 for 255.255.255.0).
For IPv6, the size value is larger due to the increased address space. IPv6 prefix lengths commonly range between /48 to /128, with /64 often used as the standard size for a single subnet.
* `address_version` - `int`: **optional**, The version of IP Address. Choose 4 for IPV4 or 6 for IPV6. Defaults to 4.
* `rawcontainer` - `boolean`: **optional**, Set to true to pass the container parameter through to the API without prefixing.
* `type` - `string`: **optional**, The Block Type for the block If not specified, a block type of Any is assumed.
* `name` - `string`: **optional**, The name of the block.
* `block_status` - `string`: **optional**, The current status of the block. 
                  Accepted values are: Deployed, FullyAssigned, Reserved, Aggregate
* `cloud_type` - `string`: **optional**, Specify the type of Cloud Provider. Currently one of: AWS, Azure, Cisco ACI, Cisco DNA Center, CloudBolt, OpenStack, ServiceNow, VMware.
* `cloud_object_id` - `string`: **optional**, The ID of this object as it is known in the cloud environment.




### Example of a Block

This example defines a data source of type `cygnalabs_ipc_subnet` and the name "my_ipc_ds", which is configured in a Terraform file.
You can reference this resource and retrieve information about it.

```hcl
data "cygnalabs_ipc_subnet" "my_ipc_ds" {
  container= "InControl/caa"
  rawcontainer=true
  address = "10.0.0.0"
  address_version=4
  size=25
}

// accessing individual field in results
output ""my-ipc-ds" {
  value = data.cygnalabs_ipc_subnet.my_ipc_ds.address 
}

```