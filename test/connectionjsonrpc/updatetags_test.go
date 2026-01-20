package connectionjsonrpc

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"

	connjsonrpc "github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc"
)

func TestUpdateTags(t *testing.T) {
	const (
		Test_HostId = "10900"
	)

	var (
		zc  *connjsonrpc.ZabbixConnectionJsonRPC
		err error
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

	t.Run("Тест 1. Инициализация соединения и получение авторизационного токена", func(t *testing.T) {
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

	t.Run("Тест 2. Обновление в хосте параметра 'теги'", func(t *testing.T) {
		res, err := zc.UpdateHostParameterTags(
			t.Context(),
			Test_HostId,
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

		hostTags, err := zc.GetHostTags(t.Context(), Test_HostId)
		assert.NoError(t, err)
		assert.Equal(t, len(hostTags), 3)
	})
}
