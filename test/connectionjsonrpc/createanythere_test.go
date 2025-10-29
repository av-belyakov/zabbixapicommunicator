package connectionjsonrpc

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc"
	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
)

func TestCreateAnyThere(t *testing.T) {
	var (
		zc *connectionjsonrpc.ZabbixConnectionJsonRPC

		err error
	)

	if err := gotenv.Load(".env"); err != nil {
		log.Fatalln(err)
	}

	zHost := os.Getenv("GO_TESTZABBIX_HOST")
	if zHost == "" {
		log.Fatalln("environment variable 'GO_TESTZABBIX_HOST' cannot be empty")
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
			connectionjsonrpc.WithHost(zHost),
			connectionjsonrpc.WithConnectionTimeout(30),
			connectionjsonrpc.WithLogin(zUser),
			connectionjsonrpc.WithPasswd(zPasswd),
		)
		assert.NoError(t, err)

		err = zc.AuthorizationStart(t.Context())
		assert.NoError(t, err)
	})

	t.Run("Тест 1. Получить информацию по API", func(t *testing.T) {
		data, err := zc.GetAPIInfo(t.Context())
		assert.NoError(t, err)

		fmt.Println("Zabbix API informetion:", string(data))
	})

	t.Run("Тест 2. Добавить новую группу хостов", func(t *testing.T) {
		newTestGroup := "ГЦМ/ ТЕСТОВАЯ ГРУППА/DEV"
		res, data, err := zc.CreateHostGroup(t.Context(), newTestGroup, 0)
		assert.NoError(t, err)

		if data.Error.Message != "" {
			fmt.Printf("Request error, code:%d, message:'%s', data:'%s'\n", data.Error.Code, data.Error.Message, data.Error.Data)

			assert.Fail(t, "request execution error", data.Error.Message, data.Error.Data)
		}

		newTestGroupId := &connectionjsonrpc.ResponseCretaeHostGroupList{}
		err = json.Unmarshal(res, newTestGroupId)
		assert.NoError(t, err)
		assert.NotEmpty(t, newTestGroupId.Result)

	})

	t.Run("Тест 3. Добавить новые хосты в группу хостов", func(t *testing.T) {

	})

	t.Cleanup(func() {
		os.Unsetenv("GO_TESTZABBIX_HOST")
		os.Unsetenv("GO_TESTZABBIX_USER")
		os.Unsetenv("GO_TESTZABBIX_PASSWD")
	})
}
