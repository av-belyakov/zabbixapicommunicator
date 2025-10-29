package connectionjsonrpc

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"

	"github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc"
)

func TestGetData(t *testing.T) {
	var (
		zc *connectionjsonrpc.ZabbixConnectionJsonRPC

		err error
	)

	if err := gotenv.Load(".env"); err != nil {
		log.Fatalln(err)
	}

	zUser := os.Getenv("GO_TESTZABBIX_USER")
	if zUser == "" {
		log.Fatalln("environment variable 'GO_TESTZABBIX_USER' cannot be empty")
	}

	zPasswd := os.Getenv("GO_TESTZABBIX_PASSWD")
	if zPasswd == "" {
		log.Fatalln("environment variable 'GO_TESTZABBIX_PASSWD' cannot be empty")
	}

	t.Run("Тест 0. Инициализация соединения и получение авторизационного токена", func(t *testing.T) {
		zc, err = connectionjsonrpc.NewConnect(
			connectionjsonrpc.WithHost("192.168.9.45"),
			connectionjsonrpc.WithConnectionTimeout(30),
			connectionjsonrpc.WithLogin(zUser),
			connectionjsonrpc.WithPasswd(zPasswd),
		)
		assert.NoError(t, err)

		err = zc.AuthorizationStart(t.Context())
		assert.NoError(t, err)
	})

	t.Run("Тест 1. Выполнение запроса", func(t *testing.T) {
		/*
				{
			  		"jsonrpc":"2.0",
			  		"method":"host.get",
			  		"params":{
			    		"search":{"host":%s}
			  		},
			  		"id":1
				}

				data, err := zc.ActionGet(t.Context(), `{
					"search":{"host": "8030174"}
					}`)
		*/
		data, err := zc.GetHostList(t.Context())
		assert.NoError(t, err)

		res := &connectionjsonrpc.ResponseMessage{}
		err = json.Unmarshal(data, res)
		assert.NoError(t, err)

		if res.Error.Message != "" {
			fmt.Printf("Request error, code:%d, message:'%s', data:'%s'\n", res.Error.Code, res.Error.Message, res.Error.Data)

			assert.Fail(t, "request execution error", res.Error.Message, res.Error.Data)
		}

		fmt.Println("Received data:", string(data))

		assert.Greater(t, len(res.Result), 0)
	})

	//	t.Run("", func(t *testing.T) {})

	t.Cleanup(func() {
		os.Unsetenv("GO_TESTZABBIX_USER")
		os.Unsetenv("GO_TESTZABBIX_PASSWD")
	})
}
