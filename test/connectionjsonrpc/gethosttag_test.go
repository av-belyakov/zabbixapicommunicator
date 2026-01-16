package connectionjsonrpc_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"

	connjsonrpc "github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc"
)

func TestGetHostTag(t *testing.T) {
	var (
		zc                *connjsonrpc.ZabbixConnectionJsonRPC
		testHostGroupName string                               = "ТЕСТ/ ГРУППА ТЕСТИРОВАНИЯ ПОЛУЧЕНИЯ ТЕГОВ/DEV"
		testHost          connjsonrpc.CreateHostOptionsRequest = connjsonrpc.CreateHostOptionsRequest{
			Host:   "test-host-for-testing-getting-tags-DEV",
			Groups: []connjsonrpc.Group{},
			Tags: []connjsonrpc.Tag{
				{Tag: "begin-tag", Value: "any begin tag"},
			},
			Macros: []connjsonrpc.Macro{
				{
					Type:        "0",
					Macro:       "{$MACRO_0}",
					Value:       "start macro",
					Description: "this is begining macro",
				},
			},
			Interfaces: connjsonrpc.InterfaceOptionsRequest{
				IP:    "126.166.78.111",
				Port:  "1996",
				DNS:   "example.domainname.org",
				Type:  1,
				Main:  1,
				Useip: 1,
				Details: connjsonrpc.DetailsOptions{
					Version: 1,
				}}}
		testHostGroupId string
		testHostId      string

		err error
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
		zc, err = connjsonrpc.NewConnect(
			connjsonrpc.WithPort(port),
			connjsonrpc.WithHost(zHost),
			connjsonrpc.WithConnectionTimeout(30),
			connjsonrpc.WithLogin(zUser),
			connjsonrpc.WithPasswd(zPasswd),
		)
		assert.NoError(t, err)

		err = zc.AuthorizationStart(t.Context())
		assert.NoError(t, err)
	})

	t.Run("Тест 1. Создать тестовую группу хостов", func(t *testing.T) {
		res, err := zc.CreateHostGroup(t.Context(), testHostGroupName)
		assert.NoError(t, err)

		hgs, errMsg, err := connjsonrpc.NewResponseCreateHostGroup().Get(res)
		assert.NoError(t, err)

		isExist := strings.ContainsAny(errMsg.Error.Message, "already exists")
		if errMsg.Error.Message != "" && !isExist {
			//fmt.Printf("Request error, code:%d, message:'%s', data:'%s'\n", data.Error.Code, data.Error.Message, data.Error.Data)
			assert.Fail(t, "request execution error", errMsg.Error.Message, errMsg.Error.Data)
		}

		if len(hgs.Result.GroupIds) != 0 {
			testHostGroupId = hgs.Result.GroupIds[0]
		}
	})

	t.Run("Тест 2. Создать тестовый хост", func(t *testing.T) {
		testHost.Groups = append(testHost.Groups, connjsonrpc.Group{GroupId: testHostGroupId})
		res, err := zc.CreateHost(t.Context(), testHost)
		assert.NoError(t, err)

		fmt.Println("Create test host, RAW response:", string(res))

		rch, errMsg, err := connjsonrpc.NewResponseCreateHost().Get(res)
		assert.NoError(t, err)
		assert.Greater(t, len(rch.Result.HostIds), 0)

		if len(rch.Result.HostIds) > 0 {
			for _, v := range rch.Result.HostIds {
				testHostId = v
			}
		}

		if errMsg.Error.Message != "" {
			fmt.Printf(
				"Request error, code:%d, message:'%s', data:'%s'\n",
				errMsg.Error.Code,
				errMsg.Error.Message,
				errMsg.Error.Data,
			)
		}

		fmt.Printf("Test host success created, id '%s'\n", testHostId)
	})

	t.Run("Тест 2. Получить информацию по тегам тестового хоста", func(t *testing.T) {
		b, err := zc.CustomRequest(
			t.Context(),
			"host.get",
			fmt.Sprintf(`{
	        "output": ["%s"],
	        "selectTags": "extend",
	        "evaltype": 0,
			"tags": []
			}`, testHostId),
		)
		assert.NoError(t, err)

		//fmt.Printf("Get RAW information for 'Tag' host with id '%s': '%s'\n", testHostId, string(b))

		result := struct {
			Result []struct {
				HostId string `json:"hostid"`
				Tags   []struct {
					Tag   string `json:"tag"`
					Value string `json:"value"`
				} `json:"tags"`
			} `json:"result"`
		}{}
		err = json.Unmarshal(b, &result)
		assert.NoError(t, err)

		fmt.Printf("Result: '%+v'\n", result)
		assert.Greater(t, len(result.Result), 0)
	})

	t.Cleanup(func() {
		//удаляем созданный хост
		if testHostId != "" {
			res, err := zc.DeleteHost(context.Background(), testHostId)
			assert.NoError(t, err)

			dh, errMsg, err := connjsonrpc.NewResponseDeleteHost().Get(res)
			assert.NoError(t, err)
			assert.Greater(t, len(dh.Result.HostIds), 0)

			fmt.Printf("Deleted responce delete host with id '%s': '%+v'\n", testHostId, dh)
			fmt.Println("Delete delete host errMsg:", errMsg)
		}

		//удаляем созданную групп хостов
		if testHostGroupId != "" {
			res, err := zc.DeleteHostGroup(context.Background(), testHostGroupId)
			assert.NoError(t, err)

			dgh, errMsg, err := connjsonrpc.NewResponseDeleteGroupHost().Get(res)
			assert.NoError(t, err)
			assert.Greater(t, len(dgh.Result.GroupIds), 0)

			fmt.Printf("Deleted responce delete group host '%+v'\n", dgh)
			fmt.Println("Delete delete group host errMsg:", errMsg)
		}

		os.Unsetenv("GO_TESTZABBIX_HOST")
		os.Unsetenv("GO_TESTZABBIX_PORT")
		os.Unsetenv("GO_TESTZABBIX_USER")
		os.Unsetenv("GO_TESTZABBIX_PASSWD")
	})
}
