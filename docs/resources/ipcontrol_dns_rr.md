# Terraform Resource: IPC DNS RR

## Overview

The `ipcontrol_dns_rr` resource manages domains resource record in IPControl.

## Parameters

### Required Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `domain` | `string` | Domain name where resource records are to be added. |
| `owner` | `string` | The owner field of the resource record. |
| `resource_record_type` | `string` | The type of resource record being imported. To import a type not predefined in IPControl the type must be prepended with “Other”. |
| `data` | `string` | The data portion of the resource record. The format is dependent on the type specified above. |

### Optional Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `domain_type` | `string` | Domain type already defined to IPControl. If not specified, the “Default” domain type will be used. |
| `resource_record_class` | `string` | The Class of the resource record. Defaults to “IN”. |
| `ttl` | `string` | Time to live. |
| `comment` | `string` | Comment text associated with the resource record. |
| `device_rec_flag` | `string` | When true, this indicates that the resource record is bound to a device. When false, this indicates that the resource record is associated with the domain only, and not a specific device. |


## Supported Resource record Type
- A
- AAAA
- A6
- AFSDB
- APL
- AVC
- CAA
- CDNSKEY
- CDS
- CERT
- CNAME
- CSYNC
- DNAME
- DS
- EUI48
- EUI64
- GPOS  
- HINFO
- HTTPS
- ISDN
- KEY
- KX
- LOC
- MB
- MG
- MINFO
- MG
- MX
- NAPTR
- NINFO
- NS
- NSAP
- NULL
- NXT
- OPENPGPKEY
- OTHER
- PX
- RKEY
- RP
- RT
- SA
- SIG
- SINK
- SRV
- SMIMEA
- SSHFP
- SVCB
- TALINK
- TKEY
- TLSA
- TSIG
- TXT
- URI
- WKS
- X25


## Example Usage

#### Domain Resource Record IPv4
```hcl
resource "ipcontrol_dns_rr" "ipv4_rr" {
  owner                = "caa"
  domain               = "com."
  resource_record_type = "A"
  data                 = "10.0.0.25"
  comment              = "RR terraform"
  ttl                  = "12345"
}
```

#### Domain Resource Record IPv6
```hcl
resource "ipcontrol_dns_rr" "ipv6_rr" {
  owner                = "caa"
  domain               = "com."
  resource_record_type = "AAAA"
  data                 = "2001:db8:85a3::3000:82"
  comment              = "RR terraform "
  ttl                  = "12345"
}
```