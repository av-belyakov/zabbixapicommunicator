package connectionjsonrpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	connjsonrpc "github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc"
	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
)

func TestGetHostParameters(t *testing.T) {
	var (
		zc            *connjsonrpc.ZabbixConnectionJsonRPC
		inventoryName string                               = "test inventory"
		testHostGroup string                               = "ТЕСТ/ Группа тестирования получения параметров хоста/DEV"
		testHost      connjsonrpc.CreateHostOptionsRequest = connjsonrpc.CreateHostOptionsRequest{
			Host:   "test-gethostparameters-host-DEV",
			Groups: []connjsonrpc.Group{},
			Tags: []connjsonrpc.Tag{
				{Tag: "begin-tag-1", Value: "any begin tag"},
				{Tag: "begin-tag-2", Value: "any begin tag"},
			},
			Macros: []connjsonrpc.Macro{
				{
					Type:        "0",
					Macro:       "{$MACRO_1}",
					Value:       "begin-macro-1",
					Description: "this is begining macro",
				},
				{
					Type:        "0",
					Macro:       "{$MACRO_2}",
					Value:       "begin-macro-2",
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
			Inventory: connjsonrpc.HostInventory{
				Name:    inventoryName,
				Alias:   "Anybody host",
				OS:      "Linux",
				OSFull:  "Linux Mint 21.1",
				Contact: "Russia, Moscow, st. Parkovay, 45",
				Vendor:  "Some vendor",
			},
			Name:          "visible host name",
			Description:   "некоторое описание хоста",
			IpmiUsername:  "someuser",
			IpmiPassword:  "somepassword",
			IpmiPrivilege: 1,
			InventoryMode: 0,
		}
		hostId      string
		hostGroupId string

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

	t.Run("Тест 1. Создание новой группы хостов", func(t *testing.T) {
		res, err := zc.CreateHostGroup(t.Context(), testHostGroup)
		assert.NoError(t, err)

		//fmt.Printf("Add group hosts, response:'%s'\n", string(res))

		hostGroup, errMsg, err := connjsonrpc.NewResponseCreateHostGroup().Get(res)
		assert.NoError(t, err)

		isExist := strings.ContainsAny(errMsg.Error.Message, "already exists")
		if errMsg.Error.Message != "" && !isExist {
			assert.Fail(t, "request execution error", errMsg.Error.Message, errMsg.Error.Data)
		}

		assert.Greater(t, len(hostGroup.Result.GroupIds), 0)

		for _, v := range hostGroup.Result.GroupIds {
			if v != "" {
				hostGroupId = v

				break
			}
		}

		fmt.Println("Created host group with ID:", hostGroupId)
	})

	t.Run("Тест 2. Создание нового хоста", func(t *testing.T) {
		testHost.Groups = []connjsonrpc.Group{
			{
				GroupId: hostGroupId,
			},
		}

		b, err := zc.CreateHost(t.Context(), testHost)
		assert.NoError(t, err)

		//fmt.Println("Response 'create host':", string(b))

		createdHost, errMsg, err := connjsonrpc.NewResponseCreateHost().Get(b)
		assert.NoError(t, err)

		isExist := strings.ContainsAny(errMsg.Error.Message, "already exists")
		if errMsg.Error.Message != "" && !isExist {
			assert.Fail(t, "request execution error", errMsg.Error.Message, errMsg.Error.Data)
		}

		assert.Greater(t, len(createdHost.Result.HostIds), 0)

		for _, v := range createdHost.Result.HostIds {
			if v != "" {
				hostId = v

				break
			}
		}

		fmt.Println("Created host with ID:", hostId)
	})

	t.Run("Тест 3. Получение параметров хоста", func(t *testing.T) {
		t.Run("Тест 3.1. Получение параметра 'группа хостов'", func(t *testing.T) {
			hostGroups, err := zc.GetHostGroups(t.Context(), hostId)
			assert.NoError(t, err)
			assert.Equal(t, len(hostGroups), 1)
		})
		t.Run("Тест 3.2. Получение параметра 'теги'", func(t *testing.T) {
			hostTags, err := zc.GetHostTags(t.Context(), hostId)
			assert.NoError(t, err)
			assert.Equal(t, len(hostTags), 2)
		})
		t.Run("Тест 3.3. Получение параметра 'макросы'", func(t *testing.T) {
			macros, err := zc.GetHostMacros(t.Context(), hostId)
			assert.NoError(t, err)
			assert.Equal(t, len(macros), 2)
		})
		t.Run("Тест 3.4. Получение параметра 'интерфейсы'", func(t *testing.T) {
			hostInterfaces, err := zc.GetHostInterface(t.Context(), hostId)
			assert.NoError(t, err)
			assert.Equal(t, len(hostInterfaces), 1)
		})
		t.Run("Тест 3.5. Получение параметра 'инвентаризация'", func(t *testing.T) {
			hostInventory, err := zc.GetHostInventory(t.Context(), hostId)
			assert.NoError(t, err)
			assert.Equal(t, hostInventory.Result[0].Inventory.Name, inventoryName)

			//fmt.Printf("\nResponse Inventory:'%#v'\n", hostInventory)
		})
	})

	t.Cleanup(func() {
		if hostId != "" {
			_, err := zc.DeleteHost(context.Background(), hostId)
			assert.NoError(t, err)

			fmt.Println("Deleted host with id:", hostId)
		}

		if hostGroupId != "" {
			_, err := zc.DeleteHostGroup(context.Background(), hostGroupId)
			assert.NoError(t, err)

			fmt.Println("Deleted host group with id:", hostGroupId)
		}

		os.Unsetenv("GO_TESTZABBIX_HOST")
		os.Unsetenv("GO_TESTZABBIX_PORT")
		os.Unsetenv("GO_TESTZABBIX_USER")
		os.Unsetenv("GO_TESTZABBIX_PASSWD")
	})
}
