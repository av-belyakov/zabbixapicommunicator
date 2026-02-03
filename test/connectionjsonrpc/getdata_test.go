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

func TestGetData(t *testing.T) {
	var (
		f  *os.File
		zc *connjsonrpc.ZabbixConnectionJsonRPC

		err error

		information connjsonrpctest.Information = connjsonrpctest.Information{}
		nameGroups  map[string]string           = map[string]string{
			"Сайты ГЦМ/ 3.1 Критические":              "",
			"Сайты ГЦМ/ 3.2 ОГВ Российской Федерации": "",
			"Сайты ГЦМ/ 3.3 ОГВ ЦФО":                  "",
			"Сайты ГЦМ/ 3.4 СМИ":                      "",
			"Сайты ГЦМ/ 3.5 ЦИК России":               "",
			"Сайты ГЦМ/ 3.6 Предприятия":              "",
			"Сайты ГЦМ/ 3.7 Россельхознадзора":        "",
		}
	)

	if err := gotenv.Load(".env.prod"); err != nil {
		log.Fatalln(err)
	}

	zHost := os.Getenv("GO_TESTZABBIX_HOST")
	if zHost == "" {
		log.Fatalln("environment variable 'GO_TESTZABBIX_HOST' cannot be empty")
	}

	tmpZPort := os.Getenv("GO_TESTZABBIX_PORT")
	if tmpZPort == "" {
		log.Fatalln("environment variable 'GO_TESTZABBIX_PORT' cannot be empty")
	}
	zPort, err := strconv.Atoi(tmpZPort)
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

	f, err = os.Create("./prod-hosts.json")
	if err != nil {
		log.Fatalln(err)
	}

	t.Run("Тест 0. Инициализация соединения и получение авторизационного токена", func(t *testing.T) {
		zc, err = connjsonrpc.NewConnect(
			connjsonrpc.WithTLS(),
			connjsonrpc.WithInsecureSkipVerify(),
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

		data, errMsg, err := connjsonrpc.NewResponseGetHostGroupList().Get(res)
		assert.NoError(t, err)

		if errMsg.Error.Message != "" {
			fmt.Printf("Request error, code:%d, message:'%s', data:'%s'\n", errMsg.Error.Code, errMsg.Error.Message, errMsg.Error.Data)

			assert.Fail(t, "request execution error", errMsg.Error.Message, errMsg.Error.Data)
		}

		var num int = 1
		//список групп
		for _, v := range data.Result {
			if _, ok := nameGroups[v.Name]; !ok {
				continue
			}

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

		//		fmt.Printf("RAW DATA:'%+v'\n", string(res))

		data, errMsg, err := connjsonrpc.NewResponseGetHostList().Get(res)
		assert.NoError(t, err)

		if errMsg.Error.Message != "" {
			fmt.Printf("Request error, code:%d, message:'%s', data:'%s'\n", errMsg.Error.Code, errMsg.Error.Message, errMsg.Error.Data)

			assert.Fail(t, "request execution error", errMsg.Error.Message, errMsg.Error.Data)
		}

		var num int = 1
		//список хостов
		for _, v := range data.Result {
			if num <= 13 {
				fmt.Printf("STRUCT:'%+v'\n", v)

				fmt.Printf(
					"%d.\n\tHostId:'%s'\n\tHost:'%s'\n\tName:'%s'\nMacros:'%v'\n",
					num,
					v.HostId,
					v.Host,
					v.Name,
					v.Macros,
				)

			}

			information.Hosts = append(information.Hosts, connjsonrpctest.HostInfo{
				HostId: v.HostId,
				Host:   v.Host,
				Name:   v.Name,
			})

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
		os.Unsetenv("GO_TESTZABBIX_USER")
		os.Unsetenv("GO_TESTZABBIX_PASSWD")
	})
}
