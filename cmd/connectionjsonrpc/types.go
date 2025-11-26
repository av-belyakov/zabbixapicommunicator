package connectionjsonrpc

import (
	"net/http"
	"time"
)

// zabbixConnectionOptions опции соединения
type zabbixConnectionOptions func(*ZabbixConnectionJsonRPC) error

// ZabbixConnectionJsonRPC соединение по протоколу JsonRPC
type ZabbixConnectionJsonRPC struct {
	connClient        *http.Client
	connectionTimeout time.Duration
	rootCAs           []string
	url               string
	host              string
	login             string
	passwd            string
	applicationType   string
	authorizationHash string
	port              int
	isTls             bool
	isCertSkipVerify  bool
}

// ZabbixAuthorizationData результат авторизации
type ZabbixAuthorizationData struct {
	Error   map[string]any `json:"error"`
	JsonRPC string         `json:"jsonrpc"`
	Result  string         `json:"result"`
	Id      int            `json:"id"`
}

// ZabbixAuthorizationErrorMessage сообщение об ошибке
type ZabbixAuthorizationErrorMessage struct {
	Data    string `json:"data"`
	Message string `json:"message"`
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

// HostInventory данные по инвенторизации
type HostInventory struct {
	Type             string `json:"type"`
	TypeFull         string `json:"type_full"`
	Name             string `json:"name"`
	Alias            string `json:"alias"`
	OS               string `json:"os"`
	OSFull           string `json:"os_full"`
	OSShort          string `json:"os_short"`
	SerialNoA        string `json:"serialno_a"`
	SerialNoB        string `json:"serialno_b"`
	Tag              string `json:"tag"`
	AssetTag         string `json:"asset_tag"`
	MacAddressA      string `json:"macaddress_a"`
	MacAddressB      string `json:"macaddress_b"`
	Hardware         string `json:"hardware"`
	HardwareFull     string `json:"hardware_full"`
	Software         string `json:"software"`
	SoftwareFull     string `json:"software_full"`
	SoftwareAppA     string `json:"software_app_a"`
	SoftwareAppB     string `json:"software_app_b"`
	SoftwareAppC     string `json:"software_app_c"`
	SoftwareAppD     string `json:"software_app_d"`
	SoftwareAppE     string `json:"software_app_e"`
	Contact          string `json:"contact"`
	Location         string `json:"location"`
	LocationLat      string `json:"location_lat"`
	LocationLon      string `json:"location_lon"`
	Notes            string `json:"notes"`
	Chassis          string `json:"chassis"`
	Model            string `json:"model"`
	HWArch           string `json:"hw_arch"`
	Vendor           string `json:"vendor"`
	ContractNumber   string `json:"contract_number"`
	InstallerName    string `json:"installer_name"`
	DeploymentStatus string `json:"deployment_status"`
	URLA             string `json:"url_a"`
	URLB             string `json:"url_b"`
	URLC             string `json:"url_c"`
	HostNetworks     string `json:"host_networks"`
	HostNetmask      string `json:"host_netmask"`
	HostRouter       string `json:"host_router"`
	OOBIP            string `json:"oob_ip"`
	OOBNetmask       string `json:"oob_netmask"`
	OOBRouter        string `json:"oob_router"`
	DateHWPurchase   string `json:"date_hw_purchase"`
	DateHWInstall    string `json:"date_hw_install"`
	DateHWExpiry     string `json:"date_hw_expiry"`
	DateHWDecomm     string `json:"date_hw_decomm"`
	SiteAddressA     string `json:"site_address_a"`
	SiteAddressB     string `json:"site_address_b"`
	SiteAddressC     string `json:"site_address_c"`
	SiteCity         string `json:"site_city"`
	SiteState        string `json:"site_state"`
	SiteCountry      string `json:"site_country"`
	SiteZip          string `json:"site_zip"`
	SiteRack         string `json:"site_rack"`
	SiteNotes        string `json:"site_notes"`
	POC1Name         string `json:"poc_1_name"`
	POC1Email        string `json:"poc_1_email"`
	POC1PhoneA       string `json:"poc_1_phone_a"`
	POC1PhoneB       string `json:"poc_1_phone_b"`
	POC1Cell         string `json:"poc_1_cell"`
	POC1Screen       string `json:"poc_1_screen"`
	POC1Notes        string `json:"poc_1_notes"`
	POC2Name         string `json:"poc_2_name"`
	POC2Email        string `json:"poc_2_email"`
	POC2PhoneA       string `json:"poc_2_phone_a"`
	POC2PhoneB       string `json:"poc_2_phone_b"`
	POC2Cell         string `json:"poc_2_cell"`
	POC2Screen       string `json:"poc_2_screen"`
	POC2Notes        string `json:"poc_2_notes"`
}
