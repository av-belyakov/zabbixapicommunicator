package connectionjsonrpc_test

import (
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
			Interfaces: connjsonrpc.InterfaceOptions{
				IP:    "26.66.78.100",
				Port:  "5996",
				DNS:   "example.domainname.org",
				Type:  1,
				Main:  1,
				Useip: 1,
				Details: connjsonrpc.DetailsOptions{
					Version: 1,
				}}}
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

		//fmt.Println("Raw response:", string(res))

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

			fmt.Println("testHostGroupId =", testHostGroupId)
			fmt.Println("testAdditionalHostGroupId =", testAdditionalHostGroupId)

			/*
				План работы:
				1. Поправить методы обновления, в настоящее время при обновлении
				какого либо параметра параметр перезатирается новым значением,
				старое значение не сохраняется так как метод не получает имеющиеся значения.

				2. Исправить ошибку 'Incorrect arguments passed to function.' при обновлении
				параметра 'интерфейсы'.

				3. Добавить метод обновления параметра 'inventory'.
			*/

			res, err = zc.UpdateHostParameterGroup(
				t.Context(),
				testHostId,
				connjsonrpc.Groups{
					Group: []connjsonrpc.Group{
						{
							GroupId: testHostGroupId,
						},
						{
							GroupId: testAdditionalHostGroupId,
						},
					}},
			)
			assert.NoError(t, err)

			fmt.Println("Raw update host parameters 'host group'", string(res))

			_, errMsgUpdate, err := connjsonrpc.NewResponseUpdateHost().Get(res)
			if errMsgUpdate.Error.Message != "" {
				assert.Fail(t, fmt.Sprintf(
					"Request error, code:%d, message:'%s', data:'%s'\n",
					errMsg.Error.Code,
					errMsg.Error.Message,
					errMsg.Error.Data,
				))
			}
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

			fmt.Println("Raw update host parameters 'tags'", string(res))

			_, errMsg, err := connjsonrpc.NewResponseUpdateHost().Get(res)
			if errMsg.Error.Message != "" {
				assert.Fail(t, fmt.Sprintf(
					"Request error, code:%d, message:'%s', data:'%s'\n",
					errMsg.Error.Code,
					errMsg.Error.Message,
					errMsg.Error.Data,
				))
			}
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

			fmt.Println("Raw update host parameters 'macro'", string(res))

			_, errMsg, err := connjsonrpc.NewResponseUpdateHost().Get(res)
			if errMsg.Error.Message != "" {
				assert.Fail(t, fmt.Sprintf(
					"Request error, code:%d, message:'%s', data:'%s'\n",
					errMsg.Error.Code,
					errMsg.Error.Message,
					errMsg.Error.Data,
				))
			}
		})
		t.Run("Тест 4.4. Обновление в хосте параметра 'интерфейсы'", func(t *testing.T) {
			res, err := zc.UpdateHostParameterInterfaces(
				t.Context(),
				testHostId,
				connjsonrpc.Interfaces{
					Interface: []connjsonrpc.InterfaceOptions{
						{
							Type:  2,
							Main:  1,
							Useip: 1,
							IP:    "127.0.0.127",
							DNS:   "dns.example.domain-name.org",
							Port:  "10051",
							//Details: connjsonrpc.DetailsOptions{
							//	Contextname:  "any context name",
							//	SecurityName: "any security name",
							//},
						}}})
			assert.NoError(t, err)

			fmt.Println("Raw update host parameters 'interface'", string(res))

			_, errMsg, err := connjsonrpc.NewResponseUpdateHost().Get(res)

			fmt.Println("|||| errMsg:", errMsg)

			if errMsg.Error.Message != "" {
				assert.Fail(t, fmt.Sprintf(
					"Request error, code:%d, message:'%s', data:'%s'\n",
					errMsg.Error.Code,
					errMsg.Error.Message,
					errMsg.Error.Data,
				))
			}
		})
	})
	t.Cleanup(func() {
		/*
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

			//удаляем дополнительно созданную групп хостов
			if testHostGroupId != "" {
				res, err := zc.DeleteHostGroup(context.Background(), testAdditionalHostGroupId)
				assert.NoError(t, err)

				dgh, errMsg, err := connjsonrpc.NewResponseDeleteGroupHost().Get(res)
				assert.NoError(t, err)
				assert.Greater(t, len(dgh.Result.GroupIds), 0)

				fmt.Printf("Deleted responce delete additional group host '%+v'\n", dgh)
				fmt.Println("Delete delete additional group host errMsg:", errMsg)
			}
		*/

		os.Unsetenv("GO_TESTZABBIX_HOST")
		os.Unsetenv("GO_TESTZABBIX_PORT")
		os.Unsetenv("GO_TESTZABBIX_USER")
		os.Unsetenv("GO_TESTZABBIX_PASSWD")
	})
}
