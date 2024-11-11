package entities

type Bool bool

// Params interface
type Params map[string]interface{}

/* IPAM Object Base interface */
type IpamObject interface {
	ObjectType() string
	// Params() Params
}

/* Object base struct */
type ObjBase struct {
	objectType string
	// Parameters Params
}

func (obj *ObjBase) ObjectType() string {
	return obj.objectType
}

// func (obj *ObjBase) Params() Params {
// 	return obj.Parameters
// }
