package zabbixapicommunicator

type Messager interface {
	GetType() string
	GetMessage() string
	SetType(v string)
	SetMessage(v string)
}
