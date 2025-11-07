package connectionjsonrpc

// CreateHostOptions опции создания хоста
type CreateHostOptions struct {
	Tags      []Tag      `json:"tags"`
	Groups    []Group    `validate:"required" json:"groups"`
	Macros    []Macro    `json:"macros"`
	Templates []Template `json:"templates"`
	Inventory struct {
		MacaddressA string `json:"macaddress_a"`
		MacaddressB string `json:"macaddress_b"`
	} `json:"inventory"`
	Interfaces    InterfacesOptions `json:"interfaces"`
	Host          string            `validate:"required" json:"host"`
	InventoryMode int               `json:"inventory_mode"`
}

// InterfacesOptions опции интерфейса хоста
type InterfacesOptions struct {
	Details DetailsOptions `json:"details"`
	IP      string         `validate:"ip" json:"ip"`
	DNS     string         `json:"dns"`
	Port    string         `json:"port"`
	HostId  string         `json:"hostid"`
	Type    int            `validate:"oneof=1 2 3 4" json:"type"`
	Main    int            `validate:"oneof=0 1" json:"main"`
	Useip   int            `validate:"oneof=0 1" json:"useip"`
}

// DetailsOptions опции деталей интерфейса
type DetailsOptions struct {
	Authpassphrase string `json:"authpassphrase"`
	Privpassphrase string `json:"privpassphrase"`
	Community      string `json:"community"`
	Contextname    string `json:"contextname"`
	SecurityName   string `json:"securityname"`
	Authprotocol   int    `validate:"oneof=0 1 2 3 4 5" json:"authprotocol"`
	Privprotocol   int    `validate:"oneof=0 1 2 3 4 5" json:"privprotocol"`
	SecurityLevel  int    `validate:"oneof=0 1 2" json:"securitylevel"`
	Version        int    `validate:"oneof=1 2 3" json:"version"`
	Bulk           int    `validate:"oneof=0 1" json:"bulk"`
	MaxRepetitions int    `json:"max_repetitions"`
}

// Group опции групп
type Group struct {
	GroupId string `json:"groupid"`
	Name    string `json:"name,omitempty"`
}

// Tag опции тегов
type Tag struct {
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

// Template опции шаблонов
type Template struct {
	TemplateId string `json:"templateid"`
}

// Macro опции макросов
type Macro struct {
	Macro       string `json:"macro"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

// CreateHostResponse ответ на запрос создания хоста
type CreateHostResponse struct {
}
