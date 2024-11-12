package models

import "terraform-provider-ipcontrol/ipcontrol/entities"

func Subnet(subnet entities.IPCSubnet) *entities.IPCSubnet {
	res := subnet
	return &res
}
