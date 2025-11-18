package connectionjsonrpc

// ResponseError ошибка
type ResponseError struct {
	Error struct {
		Message string `json:"message"`
		Data    string `json:"data"`
		Code    int    `json:"code"`
	} `json:"error"`
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

// ResponseCreateHostGroup ответ на запрос создания группы хостов
type ResponseCreateHostGroup struct {
	Result struct {
		GroupIds []string `json:"groupids"`
	} `json:"result"`
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

// ResponseCreateHost ответ на запрос создания хоста
type ResponseCreateHost struct {
	Result struct {
		HostIds []string `json:"hostids"`
	} `json:"result"`
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

// ResponseHostGroupList ответное сообщение со списком групп
type ResponseHostGroupList struct {
	Result  []HostGroupInformation `json:"result"`
	JsonRPC string                 `json:"jsonrpc"`
	ID      int                    `json:"id"`
}

// HostGroupInformation описание информации по группам
type HostGroupInformation struct {
	UUID    string `json:"uuid"`
	Name    string `json:"name"`
	Flags   string `json:"flags"`
	GroupId string `json:"groupid"`
}

// ResponseHostList ответное сообщение со списком хостов
type ResponseHostList struct {
	Result []HostInformation `json:"result"`
}

// HostInformation информация по хостам
type HostInformation struct {
	UUID              string `json:"uuid"`
	HostId            string `json:"hostid"`
	ProxyId           string `json:"proxyId"`
	Host              string `json:"host"`
	Status            string `json:"status"`
	Name              string `json:"name"`
	Flags             string `json:"flags"`
	Readme            string `json:"readme"`
	Templateid        string `json:"templateid"`
	Description       string `json:"description"`
	TlsConnect        string `json:"tls_connect"`
	TlsAccept         string `json:"tls_accept"`
	TlsIssuer         string `json:"tls_issuer"`
	TlsSubject        string `json:"tls_subject"`
	CustomInterfaces  string `json:"custom_interfaces"`
	VendorName        string `json:"vendor_name"`
	VendorVersion     string `json:"vendor_version"`
	ProxyGroupid      string `json:"proxy_groupid"`
	MonitoredBy       string `json:"monitored_by"`
	WizardReady       string `json:"wizard_ready"`
	InventoryMode     string `json:"inventory_mode"`
	ActiveAvailable   string `json:"active_available"`
	AssignedProxyid   string `json:"assigned_proxyid"`
	IpmiAuthtype      string `json:"ipmi_authtype"`
	IpmiPrivilege     string `json:"ipmi_privilege"`
	IpmiUsername      string `json:"ipmi_username"`
	IpmiPassword      string `json:"ipmi_password"`
	Maintenanceid     string `json:"maintenanceid"`
	MaintenanceStatus string `json:"maintenance_status"`
	MaintenanceType   string `json:"maintenance_type"`
	MaintenanceFrom   string `json:"maintenance_from"`
}

// ResponseUpdateHostGroup ответное сообщение на запрос обновления группы хостов
type ResponseUpdateHostGroup struct {
	//GroupIds []string `json:"groupids"`
	Result []struct {
		GroupIds string `json:"groupids"`
	} `json:"result"`
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

// ResponseUpdateHost ответное сообщение на запрос обновления хоста
type ResponseUpdateHost struct {
	Result struct {
		HostIds []string `json:"hostids"`
	} `json:"result"`
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

// ResponseDeleteGroupHost ответное сообщение на запрос удаления группы хостов
type ResponseDeleteGroupHost struct {
	Result struct {
		GroupIds []int `json:"groupids"`
	} `json:"result"`
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

// ResponseDeleteHost ответное сообщение на запрос удаления хостов
type ResponseDeleteHost struct {
	Result struct {
		HostIds []string `json:"hostids"`
	} `json:"result"`
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

// ResponseGetInterfaces оответное сообщение на запрос интерфейсов
type ResponseGetInterfaces struct {
	Result []struct {
		HostId     string              `json:"hostid"`
		Interfaces []ResponseInterface `json:"interfaces"`
	} `json:"result"`
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

type ResponseInterface struct {
	Details      DetailsOptions `json:"details"`
	Interfaceid  string         `json:"interfaceid"`
	HostId       string         `json:"hostid"`
	Useip        string         `json:"useip"`
	Type         string         `json:"type"`
	Main         string         `json:"main"`
	IP           string         `json:"ip"`
	DNS          string         `json:"dns"`
	Port         string         `json:"port"`
	Error        string         `json:"error"`
	Available    string         `json:"available"`
	ErrorsFrom   string         `json:"errors_from"`
	DisableUntil string         `json:"disable_until"`
}
