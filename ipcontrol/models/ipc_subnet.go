package models

import "terraform-provider-ipcontrol/ipcontrol/entities"

func Subnet(subnet entities.IPC_Subnet) *entities.IPC_Subnet {
	res := subnet
	return &res
}
