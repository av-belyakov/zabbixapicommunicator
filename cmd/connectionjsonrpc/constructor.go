package connectionjsonrpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"
)

func NewConnect(opts ...zabbixConnectionOptions) (*ZabbixConnectionJsonRPC, error) {
	api := &ZabbixConnectionJsonRPC{
		connectionTimeout: (1 * time.Second),
		applicationType:   "application/json-rpc",
	}

	for _, opt := range opts {
		if err := opt(api); err != nil {
			return api, err
		}
	}

	api.url = fmt.Sprintf("https://%s/api_jsonrpc.php", api.host)
	api.connClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     api.connectionTimeout,
			MaxIdleConnsPerHost: 10,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, //пока пропускаем верификацию по сертификату
				RootCAs:            x509.NewCertPool(),
			},
		},
		Timeout: 15 * time.Second,
	}

	return api, nil
}
