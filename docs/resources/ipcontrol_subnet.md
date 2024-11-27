# Resource: [IPC Subnet]

###  Descriptions
The `cygnalabs_ipc_subnet` resource associates a block with container.

### Parameters
The following list describes the parameters you can define in the resource block of the record:

* `container` - `string`: **required**, The name of the container that will hold the block.
* `address` - `string`: **required**, The address block to allocate.
* `size` - `int`: **required**, The size parameter represents the subnet mask or prefix length of the address block in CIDR notation. For IPv4, this is typically a value between 0 and 32 (e.g., 24 for 255.255.255.0).
For IPv6, the size value is larger due to the increased address space. IPv6 prefix lengths commonly range between /48 to /128, with /64 often used as the standard size for a single subnet.
* `address_version` - `int`: **optional**, The version of IP Address. Choose 4 for IPV4 or 6 for IPV6. Defaults to 4.
* `rawcontainer` - `boolean`: **optional**, Set to true to pass the container parameter through to the API without prefixing.
* `type` - `string`: **optional**, The Block Type for the block If not specified, a block type of Any is assumed.
* `dns_domain` - `string`: **optional**, The name of the dns domain that will hold the block.
* `name` - `string`: **optional**, The name of the block.
* `block_status` - `string`: **optional**, The current status of the block. 
                  Accepted values are: Deployed, FullyAssigned, Reserved, Aggregate
* `cloud_type` - `string`: **optional**, Specify the type of Cloud Provider. Currently one of: AWS, Azure, Cisco ACI, Cisco DNA Center, CloudBolt, OpenStack, ServiceNow, VMware.
* `cloud_object_id` - `string`: **optional**, The ID of this object as it is known in the cloud environment.

### ⚠️ Force Replacement Fields
The following fields after changes will require deleting and recreating the resource:
* `container` - Can't change after created in IPControl.
* `address` - Can't change after created in IPControl.
* `size` - Can't change after created in IPControl.

> **WARNING**: Changing the above fields will result in the current resource being deleted and a new one created. Make sure you back up your data and understand the impact before making changes.

## How to use
First define `resource` in the .tf file.<br>
`IPv4` example
```hcl
resource "cygnalabs_ipc_subnet" "my-subnet" {
  // required parameters
  container       = "InControl/caa"
  address         = "10.0.0.0"
  size            = 24

  // optional parameters
  
  rawcontainer    = true
  address_version = 4
  name            = "my-sunet-tf-caa"
  dns_domain      = "com"
  cloud_type      = "AWS"
  cloud_object_id = "subnet-78910"
}
```
`IPv6` expamle
```hcl
resource "cygnalabs_ipc_subnet" "my-ipc-subnet-3" {
  container = "InControl/caa"
  address = "2001:db8:85a3::3000:80"
  address_version = 6
  rawcontainer    = true
  size = 121
  name = "tf-v6-test"
}
```

Then run
```bash
terraform apply
```