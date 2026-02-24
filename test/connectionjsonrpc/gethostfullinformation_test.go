package connectionjsonrpc_test

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

func TestGetHostFullInformation(t *testing.T) {
	var (
		zc *connjsonrpc.ZabbixConnectionJsonRPC

		err error
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
			connjsonrpc.WithLogin(zUser),
			connjsonrpc.WithPasswd(zPasswd),
			connjsonrpc.WithConnectionTimeout(30),
		)
		assert.NoError(t, err)

		err = zc.AuthorizationStart(t.Context())
		assert.NoError(t, err)
	})

	t.Run("Тест 2. Запрос полной информации о существующем хосте", func(t *testing.T) {
		res, err := zc.GetFullInformationAboutHost(t.Context(), "11762")
		assert.NoError(t, err)
		assert.NotEmpty(t, res.Result)

		fmt.Printf("\n--- Result:'%#v'\n", res)
	})

	t.Run("Тест 3. Запрос полной информации о не существующем хосте", func(t *testing.T) {
		res, err := zc.GetFullInformationAboutHost(t.Context(), "98765")
		assert.Error(t, err)
		assert.Empty(t, res.Result)
	})

	t.Run("Тест 4. Запрос макросов хоста", func(t *testing.T) {
		res, err := zc.GetHostMacros(t.Context(), "11762")
		assert.NoError(t, err)
		assert.NotEmpty(t, res)

		fmt.Println("Macros:", res)
	})

	t.Run("Тест 5. Запрос тегов хоста", func(t *testing.T) {
		res, err := zc.GetHostTags(t.Context(), "11762")
		assert.NoError(t, err)
		assert.NotEmpty(t, res)

		fmt.Println("Tegs:", res)
	})

	t.Run("Тест 6. Запрос интерфейсов хоста", func(t *testing.T) {
		res, err := zc.GetHostInterface(t.Context(), "11762")
		assert.NoError(t, err)
		assert.NotEmpty(t, res)

		fmt.Println("Interfaces:", res)
	})

	t.Run("Тест 7. Завершение сеанса авторизации", func(t *testing.T) {
		ok, err := zc.Logout(t.Context())
		assert.NoError(t, err)
		assert.True(t, ok)
	})
}
