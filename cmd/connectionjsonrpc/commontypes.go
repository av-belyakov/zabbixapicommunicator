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
