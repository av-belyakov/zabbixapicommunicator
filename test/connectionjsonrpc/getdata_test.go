package connectionjsonrpc

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"

	"github.com/av-belyakov/zabbixapicommunicator/v2/cmd/connectionjsonrpc"
)

func responseDecode(data []byte) (*connectionjsonrpc.ResponseMessage, error) {
	res := &connectionjsonrpc.ResponseMessage{}
	err := json.Unmarshal(data, res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func TestGetData(t *testing.T) {
	var (
		zc *connectionjsonrpc.ZabbixConnectionJsonRPC

		err error
	)

	if err := gotenv.Load(".env"); err != nil {
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

	t.Run("Тест 0. Инициализация соединения и получение авторизационного токена", func(t *testing.T) {
		zc, err = connectionjsonrpc.NewConnect(
			connectionjsonrpc.WithHost("192.168.9.45"),
			connectionjsonrpc.WithConnectionTimeout(30),
			connectionjsonrpc.WithLogin(zUser),
			connectionjsonrpc.WithPasswd(zPasswd),
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
		data, err := zc.GetHostGroupList(t.Context())
		assert.NoError(t, err)

		res, err := responseDecode(data)
		assert.NoError(t, err)

		if res.Error.Message != "" {
			fmt.Printf("Request error, code:%d, message:'%s', data:'%s'\n", res.Error.Code, res.Error.Message, res.Error.Data)

			assert.Fail(t, "request execution error", res.Error.Message, res.Error.Data)
		}

		hostGroupList := &connectionjsonrpc.ResponseHostGroupList{}
		err = json.Unmarshal(data, hostGroupList)
		assert.NoError(t, err)
		assert.Greater(t, len(hostGroupList.Result), 0)

		//список групп
		//for k, v := range hostGroupList.Result {
		//	fmt.Printf("%d.\n\tFlags:'%s'\n\tGroupId:'%s'\n\tUUID:'%s'\n\tName:'%s'\n", k, v.Flags, v.GroupId, v.UUID, v.Name)
		//}

		//fmt.Println("List host group:", string(data))
		/*
				str := `\u0421\u0430\u0439\u0442\u044b \u0413\u0426\u041c/3.11 \u041c\u0438\u043d\u043a\u0443\u043b\u044c\u0442/\u0411\u0418\u0411\u041b\u0418\u041e\u0422\u0415\u041a\u0418`

			    decoded, err := strconv.Unquote(`"` + str + `"`)
			    if err != nil {
			        fmt.Println("Ошибка:", err)
			        return
			    }

			    fmt.Println(decoded)
		*/

		assert.Greater(t, len(res.Result), 0)

	})

	t.Run("Тест 3. Получить список всех хостов для определённой группы", func(t *testing.T) {
		/*
		   GroupId:'34'
		   Name:'Сайты ГЦМ/ 3.1 Критические'
		*/

		data, err := zc.GetHostList(t.Context(), "34")
		assert.NoError(t, err)

		res, err := responseDecode(data)
		assert.NoError(t, err)

		if res.Error.Message != "" {
			fmt.Printf("Request error, code:%d, message:'%s', data:'%s'\n", res.Error.Code, res.Error.Message, res.Error.Data)

			assert.Fail(t, "request execution error", res.Error.Message, res.Error.Data)
		}

		//fmt.Println("List host:", string(data))
		hostList := &connectionjsonrpc.ResponseHostList{}
		err = json.Unmarshal(data, hostList)
		assert.NoError(t, err)
		assert.Greater(t, len(hostList.Result), 0)

		//список хостов
		for k, v := range hostList.Result {
			fmt.Printf("%d.\n\tUUID:'%s'\n\tHostId:'%s'\n\tHost:'%s'\n\tName:'%s'\n", k, v.UUID, v.HostId, v.Host, v.Name)
		}

		/*
					Нужны ещё
					Сайты ГЦМ/ 3.1 Критические
			   		Сайты ГЦМ/ 3.2 ОГВ Российской Федерации
			   		Сайты ГЦМ/ 3.3 ОГВ ЦФО
			   		Сайты ГЦМ/ 3.4 СМИ
			   		Сайты ГЦМ/ 3.5 ЦИК России
			   		Сайты ГЦМ/ 3.6 Предприятия
			   		Сайты ГЦМ/ 3.7 Россельхознадзора
		*/

		assert.Greater(t, len(res.Result), 0)
	})

	//	t.Run("", func(t *testing.T) {})

	t.Cleanup(func() {
		os.Unsetenv("GO_TESTZABBIX_USER")
		os.Unsetenv("GO_TESTZABBIX_PASSWD")
	})
}
