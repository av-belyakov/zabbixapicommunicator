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

	"github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc"
	cjsonrpc "github.com/av-belyakov/zabbixapicommunicator/v2/test/connectionjsonrpc"
)

func TestGetAnyThereData(t *testing.T) {
	var (
		f  *os.File
		zc *connectionjsonrpc.ZabbixConnectionJsonRPC

		err error

		information cjsonrpc.Information = cjsonrpc.Information{}
		nameGroups  map[string]string    = map[string]string{}
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
		zc, err = connectionjsonrpc.NewConnect(
			connectionjsonrpc.WithHost(zHost),
			connectionjsonrpc.WithPort(zPort),
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

	t.Run("Тест 2. Получить список групп хостов", func(t *testing.T) {
		res, err := zc.GetFullHostGroupList(t.Context())
		assert.NoError(t, err)

		data, err := connectionjsonrpc.ResponseDecode(res)
		assert.NoError(t, err)

		if data.Error.Message != "" {
			fmt.Printf("Request error, code:%d, message:'%s', data:'%s'\n", data.Error.Code, data.Error.Message, data.Error.Data)

			assert.Fail(t, "request execution error", data.Error.Message, data.Error.Data)
		}

		hostGroupList := &connectionjsonrpc.ResponseHostGroupList{}
		err = json.Unmarshal(res, hostGroupList)
		assert.NoError(t, err)
		assert.Greater(t, len(hostGroupList.Result), 0)

		var num int = 1
		//список групп
		for _, v := range hostGroupList.Result {
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

		data, err := connectionjsonrpc.ResponseDecode(res)
		assert.NoError(t, err)

		if data.Error.Message != "" {
			fmt.Printf("Request error, code:%d, message:'%s', data:'%s'\n", data.Error.Code, data.Error.Message, data.Error.Data)

			assert.Fail(t, "request execution error", data.Error.Message, data.Error.Data)
		}

		//fmt.Println("List host:", string(data))
		hostList := &connectionjsonrpc.ResponseHostList{}
		err = json.Unmarshal(res, hostList)
		assert.NoError(t, err)
		assert.Greater(t, len(hostList.Result), 0)

		//список хостов
		for _, v := range hostList.Result {
			information.Hosts = append(information.Hosts, cjsonrpc.HostInfo{
				HostId: v.HostId,
				Host:   v.Host,
				Name:   v.Name,
			})
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
