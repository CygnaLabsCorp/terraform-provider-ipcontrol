package utils

type ObjectManager struct {
	connector CAAConnector
}

func NewObjectManager(connector CAAConnector) *ObjectManager {
	objMgr := new(ObjectManager)
	objMgr.connector = connector
	return objMgr
}
