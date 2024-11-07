package models

import "terraform-provider-ipcontrol/ipcontrol/entities"

func Subnet(subnet entities.Subnet) *entities.Subnet {
	res := subnet
	return &res
}
