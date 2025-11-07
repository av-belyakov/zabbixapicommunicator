package connectionjsonrpc

type Information struct {
	Hosts []HostInfo `json:"hosts"`
}

type HostInfo struct {
	HostId string `json:"host_id"`
	Host   string `json:"host"`
	Name   string `json:"name"`
}
