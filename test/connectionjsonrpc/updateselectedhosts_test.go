package connectionjsonrpc_test

import (
	"fmt"
	"log"
	"maps"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"

	"github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc"
	connjsonrpc "github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc"
	connjsonrpctest "github.com/av-belyakov/zabbixapicommunicator/v2/test/connectionjsonrpc"
)

// ---------------------------------
//будет выполнятся корректно только после вызова TestUpdateAnyThere
// ---------------------------------

func TestUpdateSelectedHosts(t *testing.T) {
	var (
		zc *connjsonrpc.ZabbixConnectionJsonRPC

		err error

		listGroupsId []string
		nameGroups   map[string]string = map[string]string{
			"ТЕСТ/ Группа тестирования получения параметров хоста/DEV": "",
			"ТЕСТ/ ДОПОЛНИТЕЛЬНАЯ ТЕСТОВАЯ ГРУППА/DEV":                 "",
			"ТЕСТ/ ТЕСТОВАЯ ГРУППА/DEV":                                "",
			"ТЕСТ/ ТЕСТОВАЯ ГРУППА/DEV_updated_name":                   "",
		}
		information connjsonrpctest.Information = connjsonrpctest.Information{}
	)

	if err := gotenv.Load(".env.test"); err != nil {
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

	t.Run("Тест 1. Инициализация соединения и получение авторизационного токена", func(t *testing.T) {
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

	t.Run("Тест 2. Получить список групп хостов", func(t *testing.T) {
		res, err := zc.GetFullHostGroupList(t.Context())
		assert.NoError(t, err)

		data, errMsg, err := connjsonrpc.NewResponseGetHostGroupList().Get(res)
		assert.NoError(t, err)

		if errMsg.Error.Message != "" {
			fmt.Printf(
				"Request error, code:%d, message:'%s', data:'%s'\n",
				errMsg.Error.Code,
				errMsg.Error.Message,
				errMsg.Error.Data,
			)
		}

		var num int = 1
		//список групп
		for _, v := range data.Result {
			if _, ok := nameGroups[v.Name]; ok {
				nameGroups[v.Name] = v.GroupId
			}

			fmt.Printf("%d.\n\tGroupId:'%s'\n\tUUID:'%s'\n\tName:'%s'\n", num, v.GroupId, v.UUID, v.Name)
			num++
		}
		assert.Greater(t, len(data.Result), 0)
	})

	t.Run("Тест 3. Получить список всех хостов для определённой группы", func(t *testing.T) {
		for v := range maps.Values(nameGroups) {
			listGroupsId = append(listGroupsId, v)
		}

		fmt.Println("listGroupsId:", listGroupsId)

		res, err := zc.GetHostList(t.Context(), listGroupsId...)
		assert.NoError(t, err)

		data, errMsg, err := connjsonrpc.NewResponseGetHostList().Get(res)
		assert.NoError(t, err)

		if errMsg.Error.Message != "" {
			fmt.Printf("Request error, code:%d, message:'%s', data:'%s'\n", errMsg.Error.Code, errMsg.Error.Message, errMsg.Error.Data)

			assert.Fail(t, "request execution error", errMsg.Error.Message, errMsg.Error.Data)
		}
		assert.Greater(t, len(data.Result), 0)

		fmt.Println("Total host groups:", len(data.Result), ", print first 10 items...")
		var num int = 1
		//список хостов
		for _, v := range data.Result {
			information.Hosts = append(information.Hosts, connjsonrpctest.HostInfo{
				HostId: v.HostId,
				Host:   v.Host,
				Name:   v.Name,
			})

			if num <= 10 {
				fmt.Printf("%d.\n\tHostId:'%s'\n\tHost:'%s'\n\tName:'%s'\n", num, v.HostId, v.Host, v.Name)

				continue
			}

			num++
		}

		assert.Greater(t, len(data.Result), 0)
	})

	t.Run("Тест 4. Заменить информацию по тегам тестового хоста", func(t *testing.T) {
		hostId := "11616" //"10897"

		zc.UpdateHostParameterTags(
			t.Context(),
			hostId,
			connjsonrpc.Tags{
				Tag: []connjsonrpc.Tag{
					{
						Tag:   "begin-tag-1",
						Value: "1236vvv",
					},
					{
						Tag:   "begin-tag-3",
						Value: "0000",
					},
				},
			},
		)

		tagList, err := zc.GetHostTags(t.Context(), hostId)
		assert.NoError(t, err)
		assert.NotEmpty(t, tagList)

		fmt.Println("tagList:", tagList)
	})

	t.Run("Тест 5. Получить информацию по тегам для тестового хоста с определённым номером", func(t *testing.T) {
		hostId := "10957"

		tags, err := zc.GetHostTags(
			t.Context(),
			hostId,
		)

		//fmt.Println("___ ERROR:", err, "____")

		assert.NoError(t, err)
		assert.NotEmpty(t, tags)

		fmt.Println("tagList:", tags)
	})

	t.Run("Тест 6. Получить список всех хостов", func(t *testing.T) {
		hostIds := []string(nil)

		res, err := zc.GetHosts(t.Context())
		assert.NoError(t, err)

		hostList, errMsg, err := connectionjsonrpc.NewResponseGetHostList().Get(res)
		assert.NoError(t, err)

		if errMsg.Error.Message != "" {
			fmt.Printf("Request error, code:%d, message:'%s', data:'%s'\n", errMsg.Error.Code, errMsg.Error.Message, errMsg.Error.Data)

			assert.Fail(t, "request execution error", errMsg.Error.Message, errMsg.Error.Data)
		}
		assert.Greater(t, len(hostList.Result), 0)

		fmt.Println("Total hosts:", len(hostList.Result), ", print first 20 items...")
		var num int = 1
		for _, v := range hostList.Result {
			hostIds = append(hostIds, v.HostId)

			if num <= 10 {
				fmt.Printf("%d.\n\tHostId:'%s'\n\tHost:'%s'\n\tName:'%s'\n", num, v.HostId, v.Host, v.Name)

				num++
			}
		}

		fmt.Println("hosts id:", hostIds)
	})
}
