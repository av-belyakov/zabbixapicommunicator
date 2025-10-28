package connectionjsonrpc

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// New создает объект соединения с Zabbix API
func New(settings SettingsZabbixConnectionJsonRPC) (*ZabbixConnectionJsonRPC, error) {
	var zc *ZabbixConnectionJsonRPC

	connTimeout := 30 * time.Second
	if settings.ConnectionTimeout > (1 * time.Second) {
		connTimeout = settings.ConnectionTimeout
	}

	if settings.Host == "" {
		return zc, errors.New("the value 'host' should not be empty")
	}

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     connTimeout,
			MaxIdleConnsPerHost: 10,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				RootCAs:            x509.NewCertPool(),
			},
		},
		Timeout: 15 * time.Second,
	}

	return &ZabbixConnectionJsonRPC{
		url:             fmt.Sprintf("https://%s/api_jsonrpc.php", settings.Host),
		host:            settings.Host,
		login:           settings.Login,
		passwd:          settings.Passwd,
		applicationType: "application/json-rpc",
		connClient:      client,
	}, nil
}
