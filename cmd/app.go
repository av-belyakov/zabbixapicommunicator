package zabbixapicommunicator

import (
	"fmt"
	"time"
)

// New конструктор создающий обработчик соединения с API Zabbix
func New(settings SettingsZabbixConnection) (*ZabbixConnection, error) {
	var zc ZabbixConnection

	if settings.Host == "" {
		return &zc, fmt.Errorf("the value 'Host' should not be empty")
	}

	if settings.Port == 0 {
		return &zc, fmt.Errorf("the value 'Port' should not be equal '0'")
	}

	if settings.ZabbixHost == "" {
		return &zc, fmt.Errorf("the value 'ZabbixHost' should not be empty")
	}

	if settings.NetProto != "tcp" && settings.NetProto != "udp" {
		settings.NetProto = "tcp"
	}

	if settings.ConnectionTimeout == nil {
		t := time.Duration(5 * time.Second)
		settings.ConnectionTimeout = &t
	}

	zc = ZabbixConnection{
		host:        settings.Host,
		port:        settings.Port,
		netProto:    settings.NetProto,
		zabbixHost:  settings.ZabbixHost,
		connTimeout: settings.ConnectionTimeout,
		chanErr:     make(chan error),
	}

	return &zc, nil
}
