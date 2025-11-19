package connectionjsonrpc

// CreateHostOptionsRequest опции создания хоста
type CreateHostOptionsRequest struct {
	Tags          []Tag                   `json:"tags"`
	Groups        []Group                 `validate:"required" json:"groups"`
	Macros        []Macro                 `json:"macros"`
	Templates     []Template              `json:"templates"`
	Inventory     HostInventory           `json:"inventory"`
	Interfaces    InterfaceOptionsRequest `json:"interfaces"`
	Name          string                  `json:"name"`
	Host          string                  `validate:"required" json:"host"`
	Description   string                  `json:"description"`
	IpmiUsername  string                  `json:"ipmi_username"`
	IpmiPassword  string                  `json:"ipmi_password"`
	InventoryMode int                     `validate:"oneof=-1 0 1" json:"inventory_mode"`          // -1 - не использовать, 0 - в ручную, 1 - автоматически
	IpmiAuthtype  int                     `validate:"oneof=-1 0 1 2 3 4 5 6" json:"ipmi_authtype"` // IPMI алгоритм аутентификации, -1 - default, 0 - none, 1 - MD2, 2 - MD5, 4 - straight, 5 - OEM, 6 - RMCP+
	IpmiPrivilege int                     `validate:"oneof=1 2 3 4 5" json:"ipmi_privilege"`       // уровень привилегий IPMI, 1 - callback, 2 - user, 3 - operator, 4 - administrator, 5 - OEM
	Flags         int                     `validate:"oneof=0 4" json:"flags"`                      // происхождение носителя, 0 - обычный хост, 4 - хост полученный из носителя
}

// InterfacesRequest интерфейсы, для запроса
type InterfacesRequest struct {
	Interface []InterfaceOptionsRequest `json:"interface"`
}

// InterfaceOptionsRequest опции интерфейса хоста, для запроса
type InterfaceOptionsRequest struct {
	Details DetailsOptions `json:"details"`
	IP      string         `validate:"ip" json:"ip"` // обязательно для заполнения если Useip = 1
	DNS     string         `json:"dns"`              // обязательно для заполнения если Useip = 0
	Port    string         `json:"port"`
	HostId  string         `json:"hostid"`
	Type    int            `validate:"oneof=1 2 3 4" json:"type"` // 1 - agent, 2 - SNMP, 3 - IPMI, 4 - JMX
	Main    int            `validate:"oneof=0 1" json:"main"`     // 1 - основной интерфейс (может быть только один основной интерфейс)
	Useip   int            `validate:"oneof=0 1" json:"useip"`    // 1 - использовать IP-адрес, 0 - использовать DNS
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
	Flags   string `json:"flags,omitzero"`
	Name    string `json:"name,omitempty"`
	UUID    string `json:"uuid,omitempty"`
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
