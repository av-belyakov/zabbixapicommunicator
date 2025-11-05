package connectionjsonrpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ActionGet получение информации по действиям.
// Содержимое 'param' должно являтся строковым представление JSON формата с определёнными
// значениями. Подробнее о типах значений и их структуре можно узнать из официальной
// документации https://www.zabbix.com/documentation/current/en/manual/api/reference/action/get.
func (api *ZabbixConnectionJsonRPC) ActionGet(ctx context.Context, param string) ([]byte, *ResponseMessage, error) {
	return api.sendRequest(
		ctx,
		strings.NewReader(
			fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"action.get",
	  			"params": %s,
	  			"id":1
			}`, param)))
}

// ActionCreate добавления новых действий над уже имеющимися данными.
// Например, новых триггеров, 'действий обнаружения', 'действий авторегистрации' и т.д.
// Содержимое 'param' должно являтся строковым представление JSON формата с определёнными
// значениями. Подробнее о типах значений и их структуре можно узнать из официальной
// документации https://www.zabbix.com/documentation/current/en/manual/api/reference/action/create
func (api *ZabbixConnectionJsonRPC) ActionCreate(ctx context.Context, param string) ([]byte, *ResponseMessage, error) {
	return api.sendRequest(
		ctx,
		strings.NewReader(
			fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"action.create",
	  			"params": %s,
	  			"id":1
			}`, param)))
}

// ActionUpdate обновления уже существующих тригеров, 'действий обнаружения', 'действий авторегистрации' и т.д.
// Содержимое 'param' должно являтся строковым представление JSON формата с определёнными
// значениями. Подробнее о типах значений и их структуре можно узнать из официальной
// документации https://www.zabbix.com/documentation/current/en/manual/api/reference/action/update.
func (api *ZabbixConnectionJsonRPC) ActionUpdate(ctx context.Context, param string) ([]byte, *ResponseMessage, error) {
	return api.sendRequest(
		ctx,
		strings.NewReader(
			fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"action.update",
	  			"params": %s,
	  			"id":1
			}`, param)))
}

// ActionDelete удаление действий установленных ранее.
// Подробнее о типах значений и их структуре можно узнать из официальной
// документации https://www.zabbix.com/documentation/current/en/manual/api/reference/action/delete.
// В сигнатуре функции перечень id является идентификатор ранее установленных тригеров
func (api *ZabbixConnectionJsonRPC) ActionDelete(ctx context.Context, id ...string) ([]byte, *ResponseMessage, error) {
	if len(id) == 0 {
		return nil, nil, errors.New("deletion cannot be performed, the list of IDs must not be empty")
	}

	res, data, err := api.sendRequest(
		ctx,
		strings.NewReader(
			fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"action.delete",
	  			"params":[%s],
	  			"id":1
			}`, strings.Join(id, ","))))

	return res, data, err
}

// GetAPIInfo информация о версии API (запрос должен выполнятся БЕЗ авторизации)
func (api *ZabbixConnectionJsonRPC) GetAPIInfo(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", api.url, strings.NewReader(`{
    		"jsonrpc": "2.0",
    		"method": "apiinfo.version",
    		"params": [],
    		"id": 1
		}`))
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("Content-Type", "application/json-rpc")

	res, err := api.connClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("error sending the request, response status is %s", res.Status)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}

// CustomRequest позволяет создавать гибкие пользовательские запросы
// Параметр 'method' описывает метод запроса. Например, "hostgroup.get" - получить
// группу хостов или один из группы хостов.
// Параметр 'param' представляет собой JSON в стороковом виде содержащий различный
// набор параметров подобных search, filter и т.д. Подробнее о формировании настраиваемого
// запроса можно узнать из официальной документации https://www.zabbix.com/documentation/current/en/manual/api/reference
func (api *ZabbixConnectionJsonRPC) CustomRequest(ctx context.Context, method string, param string) ([]byte, *ResponseMessage, error) {
	return api.sendRequest(
		ctx,
		strings.NewReader(fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"%s",
	  			"params": %s,
	  			"id":111
			}`, method, param)))
}

// GetFullHostList весь список хостов
func (api *ZabbixConnectionJsonRPC) GetFullHostList(ctx context.Context) ([]byte, *ResponseMessage, error) {
	return api.sendRequest(
		ctx,
		strings.NewReader(`{
	  			"jsonrpc":"2.0",
	  			"method":"host.get",
	  			"params": {
					"output": "extend"
				},
	  			"id":1
			}`))
}

// GetHostLis список хостов для определённых групп
func (api *ZabbixConnectionJsonRPC) GetHostList(ctx context.Context, groupId ...string) ([]byte, *ResponseMessage, error) {
	if len(groupId) == 0 {
		return nil, nil, errors.New("value 'groupId' is not be empty")
	}

	return api.sendRequest(
		ctx,
		strings.NewReader(fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"host.get",
	  			"params": {
					"output": "extend",
					"groupids": [%s]
				},
	  			"id":1
			}`, strings.Join(groupId, ","))))
}

// GetFullHostGroupList весь список групп хостов
func (api *ZabbixConnectionJsonRPC) GetFullHostGroupList(ctx context.Context) ([]byte, *ResponseMessage, error) {
	return api.sendRequest(
		ctx,
		strings.NewReader(`{
	  			"jsonrpc":"2.0",
	  			"method":"hostgroup.get",
	  			"params": {
					"output": "extend"
				},
	  			"id":1
			}`))
}

// CreateHostGroup создание группы хостов (обязательны права супер-администратора)
func (api *ZabbixConnectionJsonRPC) CreateHostGroup(ctx context.Context, name string) ([]byte, *ResponseMessage, error) {
	res, data, err := api.sendRequest(
		ctx,
		strings.NewReader(fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"hostgroup.create",
	  			"params": {
					"name":"%s" 
				},
	  			"id":1
			}`, name)))
	if err != nil {
		return nil, data, err
	}

	return res, data, err
}

// CreateHost создание хоста (обязательны права супер-администратора)
func (api *ZabbixConnectionJsonRPC) CreateHost(ctx context.Context, host string) ([]byte, *ResponseMessage, error) {
	type HostInterfaces struct {
		HostIP    string `json:"ip"`
		HostDNS   string `json:"dns"`
		HostPort  string `json:"port"`
		HostType  int    `jsom:"type"`
		HostMain  int    `json:"main"`
		HostUseip int    `json:"useip"`
	}

	res, data, err := api.sendRequest(
		ctx,
		strings.NewReader(fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"host.create",
	  			"params": {
					"host": "Linux server",
        			"interfaces": [
            {
                "type": 1,
                "main": 1,
                "useip": 1,
                "ip": "192.168.3.1",
                "dns": "",
                "port": "10050"
            }
        ],
        			"groups": [
            {
                "groupid": "50"
            }
        ],
        			"tags": [
            {
                "tag": "host-name",
                "value": "linux-server"
            }
        ],
        		"templates": [
            {
                "templateid": "20045"
            }
        ],
        		"macros": [
            {
                "macro": "{$USER_ID}",
                "value": "123321"
            },
            {
                "macro": "{$USER_LOCATION}",
                "value": "0:0:0",
                "description": "latitude, longitude and altitude coordinates"
            }
        ],
        		"inventory_mode": 0,
        		"inventory": {
            "macaddress_a": "01234",
            "macaddress_b": "56768"
        }
				},
	  			"id":1
			}`, name)))
	if err != nil {
		return nil, data, err
	}

	return res, data, err
}
