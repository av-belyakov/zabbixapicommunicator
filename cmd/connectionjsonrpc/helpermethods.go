package connectionjsonrpc

import (
	"errors"
	"time"
)

//******************* настройка опций пакета ***********************

// WithHost имя или ip адрес хоста API
func WithHost(v string) zabbixConnectionOptions {
	return func(api *ZabbixConnectionJsonRPC) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		api.host = v

		return nil
	}
}

// WithLogin имя пользователя
func WithLogin(v string) zabbixConnectionOptions {
	return func(api *ZabbixConnectionJsonRPC) error {
		if v == "" {
			return errors.New("the value of 'login' cannot be empty")
		}

		api.login = v

		return nil
	}
}

// WithPasswd пароль пользователя
func WithPasswd(v string) zabbixConnectionOptions {
	return func(api *ZabbixConnectionJsonRPC) error {
		if v == "" {
			return errors.New("the value of 'password' cannot be empty")
		}

		api.passwd = v

		return nil
	}
}

// WithConnectionTimeout временной интервал соединения в секундах
func WithConnectionTimeout(v int) zabbixConnectionOptions {
	return func(n *ZabbixConnectionJsonRPC) error {
		if v <= 1 || v > 180 {
			return errors.New("an incorrect value, the value should be in the range from 1 to 1800")
		}

		n.connectionTimeout = time.Duration(v) * time.Second

		return nil
	}
}

/*
// WithPort сетевой порт API
func WithPort(v int) zabbixConnectionOptions {
	return func(n *ZabbixConnectionJsonRPC) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		n.port = v

		return nil
	}
}
*/
