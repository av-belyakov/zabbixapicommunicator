package connectionjsonrpc_test

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"

	connjsonrpc "github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc"
	connjsonrpctest "github.com/av-belyakov/zabbixapicommunicator/v2/test/connectionjsonrpc"
)

func TestGetAnyThereData(t *testing.T) {
	var (
		f  *os.File
		zc *connjsonrpc.ZabbixConnectionJsonRPC

		err error

		information connjsonrpctest.Information = connjsonrpctest.Information{}
		nameGroups  map[string]string           = map[string]string{}
	)

	if err := gotenv.Load(".env"); err != nil {
		log.Fatalln(err)
	}

	zHost := os.Getenv("GO_TESTZABBIX_HOST")
	if zHost == "" {
		log.Fatalln("environment variable 'GO_TESTZABBIX_HOST' cannot be empty")
	}

	tmpPort := os.Getenv("GO_TESTZABBIX_PORT")
	if tmpPort == "" {
		log.Fatalln("environment variable 'GO_TESTZABBIX_PORT' cannot be empty")
	}
	zPort, err := strconv.Atoi(tmpPort)
	if err != nil {
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

	f, err = os.Create("./hosts.json")
	if err != nil {
		log.Fatalln(err)
	}

	t.Run("Тест 0. Инициализация соединения и получение авторизационного токена", func(t *testing.T) {
		zc, err = connjsonrpc.NewConnect(
			connjsonrpc.WithHost(zHost),
			connjsonrpc.WithPort(zPort),
			connjsonrpc.WithConnectionTimeout(30),
			connjsonrpc.WithLogin(zUser),
			connjsonrpc.WithPasswd(zPasswd),
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

	t.Run("Тест 2. Получить список групп хостов", func(t *testing.T) {
		res, err := zc.GetFullHostGroupList(t.Context())
		assert.NoError(t, err)

		rchg := connjsonrpc.NewResponseGetHostGroupList()
		data, errMsg, err := rchg.Get(res)
		assert.NoError(t, err)

		if errMsg.Error.Message != "" {
			fmt.Printf(
				"Request error, code:%d, message:'%s', data:'%s'\n",
				errMsg.Error.Code,
				errMsg.Error.Message,
				errMsg.Error.Data,
			)
		}
		assert.Greater(t, len(data.Result), 0)

		var num int = 1
		//список групп
		for _, v := range data.Result {
			nameGroups[v.Name] = v.GroupId

			fmt.Printf("%d.\n\tGroupId:'%s'\n\tUUID:'%s'\n\tName:'%s'\n", num, v.GroupId, v.UUID, v.Name)
			num++
		}
		assert.Greater(t, len(data.Result), 0)
	})

	t.Run("Тест 3. Получить список всех хостов для определённой группы", func(t *testing.T) {
		var listGroupsId []string
		for v := range maps.Values(nameGroups) {
			listGroupsId = append(listGroupsId, v)
		}

		res, err := zc.GetHostList(t.Context(), listGroupsId...)
		assert.NoError(t, err)

		rhl := connjsonrpc.NewResponseGetHostList()
		data, errMsg, err := rhl.Get(res)
		assert.NoError(t, err)

		if errMsg.Error.Message != "" {
			fmt.Printf("Request error, code:%d, message:'%s', data:'%s'\n", errMsg.Error.Code, errMsg.Error.Message, errMsg.Error.Data)

			assert.Fail(t, "request execution error", errMsg.Error.Message, errMsg.Error.Data)
		}
		assert.Greater(t, len(data.Result), 0)

		var num int = 1
		//список хостов
		for _, v := range data.Result {
			information.Hosts = append(information.Hosts, connjsonrpctest.HostInfo{
				HostId: v.HostId,
				Host:   v.Host,
				Name:   v.Name,
			})

			fmt.Printf("%d.\n\tHostId:'%s'\n\tHost:'%s'\n\tName:'%s'\n", num, v.HostId, v.Host, v.Name)
			num++
		}

		b, err := json.Marshal(information)
		assert.NoError(t, err)

		_, err = f.Write(b)
		assert.NoError(t, err)

		assert.Greater(t, len(data.Result), 0)
	})

	//	t.Run("", func(t *testing.T) {})

	t.Cleanup(func() {
		f.Close()

		os.Unsetenv("GO_TESTZABBIX_HOST")
		os.Unsetenv("GO_TESTZABBIX_PORT")
		os.Unsetenv("GO_TESTZABBIX_USER")
		os.Unsetenv("GO_TESTZABBIX_PASSWD")
	})
}
