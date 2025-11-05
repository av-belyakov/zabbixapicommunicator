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

	proto := "http"
	tlsConf := &tls.Config{}
	if api.isTls {
		proto = "https"

		certPool := x509.NewCertPool()
		for _, cert := range api.rootCAs {
			if cert != "" {
				certPool.AppendCertsFromPEM([]byte(cert))
			}
		}

		tlsConf = &tls.Config{
			//InsecureSkipVerify: true, //пока пропускаем верификацию по сертификату
			RootCAs: certPool,
		}
	}

	api.url = fmt.Sprintf("%s://%s:%d/api_jsonrpc.php", proto, api.host, api.port)
	api.connClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     api.connectionTimeout,
			MaxIdleConnsPerHost: 10,
			TLSClientConfig:     tlsConf,
		},
		Timeout: 15 * time.Second,
	}

	return api, nil
}
