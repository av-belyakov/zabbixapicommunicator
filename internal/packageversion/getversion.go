package packageversion

import "github.com/av-belyakov/zabbixapicommunicator/v2/constants"

// GetPackageVersion имя пакета
func GetPackageVersion() string {
	return constants.Packagr_Version
}
