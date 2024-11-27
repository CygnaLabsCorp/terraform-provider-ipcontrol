# IPControl Provider for Terraform

The following is an example contents of a provider configuration file named main.tf:

```
provider "ipcontrol" {
  server = "127.0.0.1"
  port = "1880"
  context = "workflow"
  username = "incadmin"
  password = "incadmin"
}
```

Where the fields represent the following:
- **server**: The IP address of the CAA server.
- **port**: The port used to access the CAA server.
- **context**: This is the URL root context of the CAA server. Default value is workflow.
- **username**: Username to authenticate with IPControl.
- **password**: Password to authenticate with IPControl.

## Resources

Below are the available resources for the following objectTypes:

-   IPC Subnet (ipcontrol_subnet)

## Data Sources

Below are the available IPControl data sources:

-   IPC Subnet (ipcontrol_subnet)