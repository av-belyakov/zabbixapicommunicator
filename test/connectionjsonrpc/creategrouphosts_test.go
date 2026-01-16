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
)

func TestCreateGroupHosts(t *testing.T) {
	var (
		zc *connjsonrpc.ZabbixConnectionJsonRPC

		err error

		newTestGroups []string = []string{
			"Сайты ГЦМ",
			"Сайты ГЦМ/ 3.1 Критические",
			"Сайты ГЦМ/ 3.2 ОГВ Российской Федерации",
			"Сайты ГЦМ/ 3.3 ОГВ ЦФО",
			"Сайты ГЦМ/ 3.4 СМИ",
			"Сайты ГЦМ/ 3.5 ЦИК России",
			"Сайты ГЦМ/ 3.6 Предприятия",
			"Сайты ГЦМ/ 3.7 Россельхознадзора",
		}
		newTestGroupsId map[string]string = map[string]string{}
	)

	if err := gotenv.Load(".env.test"); err != nil {
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

	t.Run("Тест 2. Добавить новую группу хостов", func(t *testing.T) {
		for _, newTestGroup := range newTestGroups {
			res, err := zc.CreateHostGroup(t.Context(), newTestGroup)
			assert.NoError(t, err)

			//fmt.Printf("Add group hosts, response:'%s'\n", string(res))

			_, errMsg, err := connjsonrpc.NewResponseCreateHostGroup().Get(res)
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

}
