package connectionjsonrpc

// CreateHostOptionsRequest опции создания хоста
type CreateHostOptionsRequest struct {
	Tags      []Tag      `json:"tags"`
	Groups    []Group    `validate:"required" json:"groups"`
	Macros    []Macro    `json:"macros"`
	Templates []Template `json:"templates"`
	Inventory struct {
		MacaddressA string `json:"macaddress_a"`
		MacaddressB string `json:"macaddress_b"`
	} `json:"inventory"`
	Interfaces    InterfaceOptions `json:"interfaces"`
	Host          string           `validate:"required" json:"host"`
	InventoryMode int              `json:"inventory_mode"`
}

// Interfaces интерфейсы
type Interfaces struct {
	Interface []InterfaceOptions `json:"interface"`
}

// InterfaceOptions опции интерфейса хоста
type InterfaceOptions struct {
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

// Groups группы
type Groups struct {
	Group []Group `json:"group"`
}

// Group опции группы
type Group struct {
	GroupId string `json:"groupid,omitempty"`
	Name    string `json:"name,omitempty"`
}

// Tags теги
type Tags struct {
	Tag []Tag `json:"tag"`
}

// Tag опции тега
type Tag struct {
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

// Template опции шаблонов
type Template struct {
	TemplateId string `json:"templateid"`
}

// Macros макросы
type Macros struct {
	Macro []Macro `json:"macro"`
}

// Macro опции макроса
type Macro struct {
	Type        string `validate:"oneof=0 1 2" json:"type"`
	Macro       string `json:"macro"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type FilterHostGroup struct {
	Name    []string `json:"name,omitempty"`
	GroupId []string `json:"groupid,omitempty"`
}
