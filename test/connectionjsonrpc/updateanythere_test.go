package connectionjsonrpc_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"

	connjsonrpc "github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc"
	responsejsonrpc "github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc/responses"
)

func TestUpdateAnyThere(t *testing.T) {
	var (
		zc                          *connjsonrpc.ZabbixConnectionJsonRPC
		testAdditionalHostGroupName string                               = "ТЕСТ/ ДОПОЛНИТЕЛЬНАЯ ТЕСТОВАЯ ГРУППА/DEV"
		testHostGroupName           string                               = "ТЕСТ/ ТЕСТОВАЯ ГРУППА/DEV"
		testHost                    connjsonrpc.CreateHostOptionsRequest = connjsonrpc.CreateHostOptionsRequest{
			Host:   "test-any-host-DEV",
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
				IP:    "26.66.78.100",
				Port:  "5996",
				DNS:   "example.domainname.org",
				Type:  1,
				Main:  1,
				Useip: 1,
				Details: connjsonrpc.DetailsOptions{
					Version: 1,
				}},
			IpmiPrivilege: 1,
		}
		testAdditionalHostGroupId string
		testHostGroupId           string
		testHostId                string

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

		hgs, errMsg, err := responsejsonrpc.NewResponseCreateHostGroup().Get(res)
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

	t.Run("Тест 2. Обновить информацию по группе хостов", func(t *testing.T) {
		res, err := zc.GetHostGroup(t.Context(), connjsonrpc.FilterHostGroup{
			GroupId: []string{testHostGroupId},
		})
		assert.NoError(t, err)

		ghl, errMsg, err := connjsonrpc.NewResponseGetHostGroupList().Get(res)
		assert.NoError(t, err)
		if len(ghl.Result) == 0 {
			fmt.Println("GetHostGroup method, error message:", errMsg)
		} else {
			assert.Equal(t, ghl.Result[0].GroupId, testHostGroupId)
		}

		newNameHostGroup := testHostGroupName + "_updated_name"
		res, err = zc.UpdateHostGroup(t.Context(), testHostGroupId, newNameHostGroup)
		assert.NoError(t, err)

		hostGroupReq, errMsg, err := connjsonrpc.NewResponseUpdateHostGroup().Get(res)
		assert.NoError(t, err)
		if errMsg.Error.Message != "" {
			fmt.Printf("Request error, code:%d, message:'%s', data:'%s'\n", errMsg.Error.Code, errMsg.Error.Message, errMsg.Error.Data)

			assert.Fail(t, "request execution error", errMsg.Error.Message, errMsg.Error.Data)
		}

		if len(hostGroupReq.Result) > 0 {
			assert.Equal(t, hostGroupReq.Result[0].GroupIds, testHostGroupId)
		}
	})

	t.Run("Тест 3. Создать тестовый хост", func(t *testing.T) {
		testHost.Groups = append(testHost.Groups, connjsonrpc.Group{GroupId: testHostGroupId})
		res, err := zc.CreateHost(t.Context(), testHost)
		assert.NoError(t, err)

		rch, errMsg, err := responsejsonrpc.NewResponseCreateHost().Get(res)
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
	})

	t.Run("Тест 4. Обновить информацию по хостам", func(t *testing.T) {
		t.Run("Тест 4.1. Обновление в хосте параметра 'группы хостов'", func(t *testing.T) {
			//создаём дополнительный тестовую группу хостов
			res, err := zc.CreateHostGroup(t.Context(), testAdditionalHostGroupName)
			assert.NoError(t, err)

			hgs, errMsg, err := responsejsonrpc.NewResponseCreateHostGroup().Get(res)
			assert.NoError(t, err)

			isExist := strings.ContainsAny(errMsg.Error.Message, "already exists")
			if errMsg.Error.Message != "" && !isExist {
				//fmt.Printf("Request error, code:%d, message:'%s', data:'%s'\n", data.Error.Code, data.Error.Message, data.Error.Data)
				assert.Fail(t, "request execution error", errMsg.Error.Message, errMsg.Error.Data)
			}

			if len(hgs.Result.GroupIds) != 0 {
				testAdditionalHostGroupId = hgs.Result.GroupIds[0]
			}

			//fmt.Println("ID testHostGroupId =", testHostGroupId)
			//fmt.Println("ID testAdditionalHostGroupId =", testAdditionalHostGroupId)

			//добавляем дополнительную группу хостов в хост
			res, err = zc.UpdateHostParameterGroups(
				t.Context(),
				testHostId,
				testAdditionalHostGroupId)
			assert.NoError(t, err)

			_, errMsgUpdate, err := connjsonrpc.NewResponseUpdateHost().Get(res)
			if errMsgUpdate.Error.Message != "" {
				assert.Fail(t, fmt.Sprintf(
					"Request error, code:%d, message:'%s', data:'%s'\n",
					errMsg.Error.Code,
					errMsg.Error.Message,
					errMsg.Error.Data,
				))
			}

			hostGroups, err := zc.GetHostGroups(t.Context(), testHostId)
			assert.NoError(t, err)
			assert.Equal(t, len(hostGroups), 2)
		})
		t.Run("Тест 4.2. Обновление в хосте параметра 'теги'", func(t *testing.T) {
			res, err := zc.UpdateHostParameterTags(
				t.Context(),
				testHostId,
				connjsonrpc.Tags{
					Tag: []connjsonrpc.Tag{
						{
							Tag:   "tag1",
							Value: "any value tag 1",
						},
						{
							Tag:   "tag2",
							Value: "any value tag 2",
						},
					}})
			assert.NoError(t, err)

			_, errMsg, err := connjsonrpc.NewResponseUpdateHost().Get(res)
			if errMsg.Error.Message != "" {
				assert.Fail(t, fmt.Sprintf(
					"Request error, code:%d, message:'%s', data:'%s'\n",
					errMsg.Error.Code,
					errMsg.Error.Message,
					errMsg.Error.Data,
				))
			}

			hostTags, err := zc.GetHostTags(t.Context(), testHostId)
			assert.NoError(t, err)
			assert.Equal(t, len(hostTags), 3)
		})
		t.Run("Тест 4.3. Обновление в хосте параметра 'макросы'", func(t *testing.T) {
			res, err := zc.UpdateHostParameterMacro(
				t.Context(),
				testHostId,
				connjsonrpc.Macros{
					Macro: []connjsonrpc.Macro{
						{
							Macro: "{$MACRO1}",
							Value: "any value macro 1",
						},
						{
							Macro:       "{$MACRO2}",
							Value:       "any value macro 2",
							Description: "any message",
						},
						{
							Macro: "{$MACRO3}",
							Value: "any value macro 3",
						}}})
			assert.NoError(t, err)

			_, errMsg, err := connjsonrpc.NewResponseUpdateHost().Get(res)
			if errMsg.Error.Message != "" {
				assert.Fail(t, fmt.Sprintf(
					"Request error, code:%d, message:'%s', data:'%s'\n",
					errMsg.Error.Code,
					errMsg.Error.Message,
					errMsg.Error.Data,
				))
			}

			macros, err := zc.GetHostMacros(t.Context(), testHostId)
			assert.NoError(t, err)
			assert.Equal(t, len(macros), 4)
		})
		t.Run("Тест 4.4. Обновление в хосте параметра 'интерфейсы'", func(t *testing.T) {
			res, err := zc.UpdateHostParameterInterfaces(
				t.Context(),
				testHostId,
				connjsonrpc.InterfacesRequest{
					Interface: []connjsonrpc.InterfaceOptionsRequest{
						{
							Type:  1,
							Main:  0,
							Useip: 0,
							IP:    "127.0.0.127",
							DNS:   "dns.example.domain-name.org",
							Port:  "10051",
							//Details: connjsonrpc.DetailsOptions{
							//	Contextname:  "any context name",
							//	SecurityName: "any security name",
							//},
						}}})
			assert.NoError(t, err)

			_, errMsg, err := connjsonrpc.NewResponseUpdateHost().Get(res)
			if errMsg.Error.Message != "" {
				assert.Fail(t, fmt.Sprintf(
					"Request error, code:%d, message:'%s', data:'%s'\n",
					errMsg.Error.Code,
					errMsg.Error.Message,
					errMsg.Error.Data,
				))
			}

			hostInterfaces, err := zc.GetHostInterface(t.Context(), testHostId)
			assert.NoError(t, err)
			assert.Equal(t, len(hostInterfaces), 2)
		})

		t.Run("Тест 4.5. Запись нового параметра 'инвентаризация', старые данные будут затёрты", func(t *testing.T) {
			inventName := "My new test inventory parameter"

			res, err := zc.CreateHostParameterInventory(
				t.Context(),
				testHostId,
				connjsonrpc.HostInventory{
					Name:     inventName,
					OS:       "MacOS",
					OSFull:   "MacOS 26.1, Tahoe",
					Type:     "simple inventory",
					TypeFull: "simple inventory for test",
					Tag:      "Host MacOS",
					Location: "Moscow",
					Contact:  "Russia, Moscow, st. Yaroslavskoe, 34",
				})
			assert.NoError(t, err)

			createdHost, errMsg, err := connjsonrpc.NewResponseCreateHost().Get(res)
			assert.NoError(t, err)

			if errMsg.Error.Message != "" {
				assert.Fail(t, fmt.Sprintf(
					"Request error, code:%d, message:'%s', data:'%s'\n",
					errMsg.Error.Code,
					errMsg.Error.Message,
					errMsg.Error.Data,
				))
			}

			assert.Greater(t, len(createdHost.Result.HostIds), 0)
			assert.Equal(t, createdHost.Result.HostIds[0], testHostId)

			hostInventory, err := zc.GetHostInventory(t.Context(), testHostId)
			assert.NoError(t, err)

			//fmt.Printf("Host inventory: %+v\n", hostInventory)

			assert.Equal(t, hostInventory.Result[0].Inventory.Name, inventName)
		})
	})
	t.Cleanup(func() {
		//удаляем созданный хост
		if testHostId != "" {
			res, err := zc.DeleteHost(context.Background(), testHostId)
			assert.NoError(t, err)

			dh, errMsg, err := connjsonrpc.NewResponseDeleteHost().Get(res)
			assert.NoError(t, err)
			assert.Greater(t, len(dh.Result.HostIds), 0)

			//fmt.Printf("Deleted responce delete host with id '%s': '%+v'\n", testHostId, dh)
			if errMsg.Error.Message != "" {
				assert.Fail(t, fmt.Sprintf(
					"Request error, code:%d, message:'%s', data:'%s'\n",
					errMsg.Error.Code,
					errMsg.Error.Message,
					errMsg.Error.Data,
				))
			}
		}

		//удаляем созданную групп хостов
		if testHostGroupId != "" {
			res, err := zc.DeleteHostGroup(context.Background(), testHostGroupId)
			assert.NoError(t, err)

			dgh, errMsg, err := connjsonrpc.NewResponseDeleteGroupHost().Get(res)
			assert.NoError(t, err)
			assert.Greater(t, len(dgh.Result.GroupIds), 0)

			//fmt.Printf("Deleted responce delete group host '%+v'\n", dgh)
			if errMsg.Error.Message != "" {
				assert.Fail(t, fmt.Sprintf(
					"Request error, code:%d, message:'%s', data:'%s'\n",
					errMsg.Error.Code,
					errMsg.Error.Message,
					errMsg.Error.Data,
				))
			}
		}

		//удаляем дополнительно созданную групп хостов
		if testHostGroupId != "" {
			res, err := zc.DeleteHostGroup(context.Background(), testAdditionalHostGroupId)
			assert.NoError(t, err)

			dgh, errMsg, err := connjsonrpc.NewResponseDeleteGroupHost().Get(res)
			assert.NoError(t, err)
			assert.Greater(t, len(dgh.Result.GroupIds), 0)

			//fmt.Printf("Deleted responce delete additional group host '%+v'\n", dgh)
			if errMsg.Error.Message != "" {
				assert.Fail(t, fmt.Sprintf(
					"Request error, code:%d, message:'%s', data:'%s'\n",
					errMsg.Error.Code,
					errMsg.Error.Message,
					errMsg.Error.Data,
				))
			}
		}

		os.Unsetenv("GO_TESTZABBIX_HOST")
		os.Unsetenv("GO_TESTZABBIX_PORT")
		os.Unsetenv("GO_TESTZABBIX_USER")
		os.Unsetenv("GO_TESTZABBIX_PASSWD")
	})
}
