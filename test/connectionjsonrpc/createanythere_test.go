package connectionjsonrpc_test

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"

	connjsonrpc "github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc"
	responsejsonrpc "github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc/responses"
)

func TestCreateAnyThere(t *testing.T) {
	var (
		zc *connjsonrpc.ZabbixConnectionJsonRPC

		err error

		newTestGroups []string = []string{
			"ГЦМ/ ТЕСТОВАЯ ГРУППА ГЦМ/DEV",
			"РЦМ/ ТЕСТОВАЯ ГРУППА РЦМ-Ставрополь/DEV",
			"РЦМ/ ТЕСТОВАЯ ГРУППА РЦМ-Симферополь/DEV",
			"РЦМ/ ТЕСТОВАЯ ГРУППА РЦМ-Москва/DEV",
			"РЦМ/ ТЕСТОВАЯ ГРУППА РЦМ-Нижний-Новгород/DEV",
			"РЦМ/ ТЕСТОВАЯ ГРУППА РЦМ-Смоленск/DEV",
			"РЦМ/ ТЕСТОВАЯ ГРУППА РЦМ-Хабаровск/DEV",
		}
		newTestGroupsId map[string]string = map[string]string{}
		newTestHosts    map[string]struct {
			Name string
			Ip   string
			DNS  string
			Port int
		} = map[string]struct {
			Name string
			Ip   string
			DNS  string
			Port int
		}{
			"test host one": {
				Name: "host one",
				Ip:   "30.122.36.76",
				DNS:  "example-one.domainname.org",
				Port: 3475,
			},
			"test host two": {
				Name: "host two",
				Ip:   "12.47.66.113",
				DNS:  "example-two.domainname.org",
				Port: 7633,
			},
			"test host three": {
				Name: "host three",
				Ip:   "91.100.32.66",
				DNS:  "example-three.domainname.org",
				Port: 9663,
			},
		}
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

	t.Run("Тест 1. Получить информацию по API", func(t *testing.T) {
		data, err := zc.GetAPIInfo(t.Context())
		assert.NoError(t, err)

		log.Println("Zabbix API informetion:", string(data))
	})

	t.Run("Тест 2. Добавить новую группу хостов", func(t *testing.T) {
		for _, newTestGroup := range newTestGroups {
			res, err := zc.CreateHostGroup(t.Context(), newTestGroup)
			assert.NoError(t, err)

			//fmt.Printf("Add group hosts, response:'%s'\n", string(res))

			rchg := responsejsonrpc.NewResponseCreateHostGroup()
			_, errMsg, err := rchg.Get(res)
			assert.NoError(t, err)

			isExist := strings.ContainsAny(errMsg.Error.Message, "already exists")

			if errMsg.Error.Message != "" && !isExist {
				//fmt.Printf("Request error, code:%d, message:'%s', data:'%s'\n", data.Error.Code, data.Error.Message, data.Error.Data)

				assert.Fail(t, "request execution error", errMsg.Error.Message, errMsg.Error.Data)
			}
		}

		//получить список групп хостов
		res, err := zc.GetFullHostGroupList(t.Context())
		assert.NoError(t, err)

		rchg := responsejsonrpc.NewResponseGetHostGroupList()
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

		for _, result := range data.Result {
			//fmt.Printf("result.name = '%s', type: %T\n", result["name"], result["name"])
			//nameGroup := fmt.Sprint(groupId["name"])

			if slices.Contains(newTestGroups, result.Name) {
				newTestGroupsId[result.GroupId] = result.Name
			}
		}
		assert.Greater(t, len(newTestGroupsId), 0)
		assert.Equal(t, len(newTestGroupsId), len(newTestGroups))
	})

	t.Run("Тест 3. Добавить новые хосты в группу хостов", func(t *testing.T) {
		var isError bool

		var groups []connjsonrpc.Group
		for groupId := range newTestGroupsId {
			//fmt.Println("___ groupId:", groupId)

			groups = append(groups, connjsonrpc.Group{GroupId: groupId})
		}

		for k, v := range newTestHosts {
			res, err := zc.CreateHost(t.Context(), connjsonrpc.CreateHostOptionsRequest{
				Host:   k,
				Groups: groups,
				Interfaces: connjsonrpc.InterfacesOptions{
					IP:    v.Ip,
					Port:  fmt.Sprint(v.Port),
					DNS:   v.DNS,
					Type:  1,
					Main:  1,
					Useip: 1,
					Details: connjsonrpc.DetailsOptions{
						Version: 1,
					},
				},
			})
			assert.NoError(t, err)
			if err != nil {
				isError = true

				break
			}

			//fmt.Println("Response:", string(res))
			rch := responsejsonrpc.NewResponseCreateHost()
			_, errMsg, err := rch.Get(res)
			assert.NoError(t, err)

			if errMsg.Error.Message != "" {
				fmt.Printf(
					"Request error, code:%d, message:'%s', data:'%s'\n",
					errMsg.Error.Code,
					errMsg.Error.Message,
					errMsg.Error.Data,
				)
			}
		}
		assert.False(t, isError)
	})

	t.Cleanup(func() {
		os.Unsetenv("GO_TESTZABBIX_HOST")
		os.Unsetenv("GO_TESTZABBIX_PORT")
		os.Unsetenv("GO_TESTZABBIX_USER")
		os.Unsetenv("GO_TESTZABBIX_PASSWD")
	})
}
