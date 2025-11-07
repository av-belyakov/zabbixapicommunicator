package connectionjsonrpc

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"

	"github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc"
)

func TestCreateAnyThere(t *testing.T) {
	var (
		zc *connectionjsonrpc.ZabbixConnectionJsonRPC

		err error

		newTestGroup   string = "ГЦМ/ ТЕСТОВАЯ ГРУППА/DEV"
		newTestGroupId string
	)

	if err := gotenv.Load(".env"); err != nil {
		log.Fatalln(err)
	}

	zHost := os.Getenv("GO_TESTZABBIX_HOST")
	if zHost == "" {
		log.Fatalln("environment variable 'GO_TESTZABBIX_HOST' cannot be empty")
	}

	zPort := os.Getenv("GO_TESTZABBIX_PORT")
	if zPort == "" {
		log.Fatalln("environment variable 'GO_TESTZABBIX_PORT' cannot be empty")
	}

	zUser := os.Getenv("GO_TESTZABBIX_USER")
	if zUser == "" {
		log.Fatalln("environment variable 'GO_TESTZABBIX_USER' cannot be empty")
	}

	zPasswd := os.Getenv("GO_TESTZABBIX_PASSWD")
	if zPasswd == "" {
		log.Fatalln("environment variable 'GO_TESTZABBIX_PASSWD' cannot be empty")
	}

	port, err := strconv.Atoi(zPort)
	if err != nil {
		log.Fatalln(err)
	}

	t.Run("Тест 0. Инициализация соединения и получение авторизационного токена", func(t *testing.T) {
		zc, err = connectionjsonrpc.NewConnect(
			connectionjsonrpc.WithPort(port),
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

		log.Println("Zabbix API informetion:", string(data))
	})

	t.Run("Тест 2. Добавить новую группу хостов", func(t *testing.T) {
		res, err := zc.CreateHostGroup(t.Context(), newTestGroup)
		assert.NoError(t, err)

		data, err := connectionjsonrpc.ResponseDecode(res)
		assert.NoError(t, err)

		isExist := strings.ContainsAny(data.Error.Message, "already exists")

		if data.Error.Message != "" && !isExist {
			//fmt.Printf("Request error, code:%d, message:'%s', data:'%s'\n", data.Error.Code, data.Error.Message, data.Error.Data)

			assert.Fail(t, "request execution error", data.Error.Message, data.Error.Data)
		}

		if !isExist {
			newTestGroupId := &connectionjsonrpc.ResponseCretaeHostGroupList{}
			err = json.Unmarshal(res, newTestGroupId)
			assert.NoError(t, err)
			assert.NotEmpty(t, newTestGroupId.Result)
		}

		res, err = zc.GetFullHostGroupList(t.Context())
		assert.NoError(t, err)

		data, err = connectionjsonrpc.ResponseDecode(res)
		assert.NoError(t, err)

		var groupIsExist bool
		for _, result := range data.Result {
			//fmt.Printf("result.name = '%s', type: %T\n", result["name"], result["name"])

			if result["name"] == newTestGroup {
				newTestGroupId = fmt.Sprint(result["groupid"])
				groupIsExist = true
			}
		}
		assert.True(t, groupIsExist)
		assert.NotNil(t, newTestGroupId)
	})

	t.Run("Тест 3. Добавить новые хосты в группу хостов", func(t *testing.T) {
		res, err := zc.CreateHost(t.Context(), connectionjsonrpc.CreateHostOptions{
			Host: "My new test host",
			Groups: []connectionjsonrpc.Group{
				{
					GroupId: newTestGroupId,
				},
			},
			Interfaces: connectionjsonrpc.InterfacesOptions{
				IP:    "45.63.22.31",
				Port:  "7899",
				DNS:   "anythere.domain.name.org",
				Type:  1,
				Main:  1,
				Useip: 1,
				Details: connectionjsonrpc.DetailsOptions{
					Version: 1,
				},
			},
		})
		assert.NoError(t, err)

		fmt.Println("Response:", string(res))
	})

	t.Cleanup(func() {
		os.Unsetenv("GO_TESTZABBIX_HOST")
		os.Unsetenv("GO_TESTZABBIX_PORT")
		os.Unsetenv("GO_TESTZABBIX_USER")
		os.Unsetenv("GO_TESTZABBIX_PASSWD")
	})
}
