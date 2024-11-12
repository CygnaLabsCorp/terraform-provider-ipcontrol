package entities

import (
	"encoding/json"
	"fmt"
	"strings"
)

type IPCSubnet struct {
	ObjBase
	Container []string `json:"container,omitempty"`
	Address   string   `json:"blockAddr,omitempty"`
	Type      string   `json:"blockType,omitempty"`
	Size      int      `json:"blockSize,omitempty"`
}

/*
 * Subnet object constructor
 */
func NewSubnet(sb IPCSubnet) *IPCSubnet {
	res := sb
	res.objectType = "subnet"
	return &res
}

func (b Bool) MarshalJSON() ([]byte, error) {
	if b {
		return json.Marshal("True")
	}
	return json.Marshal("False")
}

func (m IPCSubnet) String() string {
	return fmt.Sprintf(
		"IPC_Subnet{Container: [%s], Address: %s, Type: %s, Size: %d}",
		strings.Join(m.Container, ", "),
		m.Address,
		m.Type,
		m.Size,
	)
}
