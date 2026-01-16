package connectionjsonrpc

// ResponseError общее описание ошибок
type ResponseError struct {
	JsonRPC string `json:"jsonrpc"`
	Error   struct {
		Message string `json:"message"`
		Data    string `json:"data"`
		Code    int    `json:"code"`
	} `json:"error"`
	ID int `json:"id"`
}

// ResponseAPIInfo ответное сообщение на запрос GetAPIInfo
type ResponseAPIInfo struct {
	Result  string `json:"result"`
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
	Result  []HostInformation `json:"result"`
	JsonRPC string            `json:"jsonrpc"`
	ID      int               `json:"id"`
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
	JsonRPC string `json:"jsonrpc"`
	Result  struct {
		GroupIds []int `json:"groupids"`
	} `json:"result"`
	ID int `json:"id"`
}

// ResponseDeleteHost ответное сообщение на запрос удаления хостов
type ResponseDeleteHost struct {
	Result struct {
		HostIds []string `json:"hostids"`
	} `json:"result"`
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

// ResponseInterface описание интерфейса
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

// ResponseInventory инвенторизация хоста
type ResponseInventory struct {
	Result []struct {
		Host      string        `json:"host"`
		HostId    string        `json:"hostid"`
		Inventory HostInventory `json:"inventory"`
	} `json:"result"`
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
}
