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
		connTimeout: *settings.ConnectionTimeout,
		chanErr:     make(chan error),
	}

	return &zc, nil
}

/*
==================
WARNING: DATA RACE
Write at 0x00c000191948 by goroutine 13:
  runtime.racewrite()
      <autogenerated>:1 +0x1e
  github.com/av-belyakov/zabbixapicommunicator/cmd.(*ZabbixConnection).Start.func1()
      /home/artemij/go/pkg/mod/github.com/av-belyakov/zabbixapicommunicator@v0.0.0-20250122111136-eeef7fcc6fc4/cmd/methods.go:25 +0x33

Previous read at 0x00c000191948 by main goroutine:
  runtime.raceread()
      <autogenerated>:1 +0x1e
  github.com/av-belyakov/zabbixapicommunicator/cmd.(*ZabbixConnection).Start()
      /home/artemij/go/pkg/mod/github.com/av-belyakov/zabbixapicommunicator@v0.0.0-20250122111136-eeef7fcc6fc4/cmd/methods.go:38 +0x304
  github.com/av-belyakov/thehivehook_go_package/cmd/wrappers.WrappersZabbixInteraction()
      /home/artemij/go/src/thehivehook_go_package/cmd/wrappers/wrappers.go:46 +0x5de
  main.server()
      /home/artemij/go/src/thehivehook_go_package/cmd/server.go:80 +0xae4
  main.main()
      /home/artemij/go/src/thehivehook_go_package/cmd/main.go:26 +0x1b9

Goroutine 13 (running) created at:
  github.com/av-belyakov/zabbixapicommunicator/cmd.(*ZabbixConnection).Start()
      /home/artemij/go/pkg/mod/github.com/av-belyakov/zabbixapicommunicator@v0.0.0-20250122111136-eeef7fcc6fc4/cmd/methods.go:24 +0x178
  github.com/av-belyakov/thehivehook_go_package/cmd/wrappers.WrappersZabbixInteraction()
      /home/artemij/go/src/thehivehook_go_package/cmd/wrappers/wrappers.go:46 +0x5de
  main.server()
      /home/artemij/go/src/thehivehook_go_package/cmd/server.go:80 +0xae4
  main.main()
      /home/artemij/go/src/thehivehook_go_package/cmd/main.go:26 +0x1b9
==================
*/
