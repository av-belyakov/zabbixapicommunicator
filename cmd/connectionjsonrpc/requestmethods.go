package connectionjsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/av-belyakov/zabbixapicommunicator/v2/internal/supportingfunctions"
	"github.com/go-playground/validator/v10"
)

// ActionGet получение информации по действиям.
// Содержимое 'param' должно являтся строковым представление JSON формата с определёнными
// значениями. Подробнее о типах значений и их структуре можно узнать из официальной
// документации https://www.zabbix.com/documentation/current/en/manual/api/reference/action/get.
func (api *ZabbixConnectionJsonRPC) ActionGet(ctx context.Context, param string) ([]byte, error) {
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
func (api *ZabbixConnectionJsonRPC) ActionCreate(ctx context.Context, param string) ([]byte, error) {
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
func (api *ZabbixConnectionJsonRPC) ActionUpdate(ctx context.Context, param string) ([]byte, error) {
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

// ActionDelete удаление действий.
// Подробнее о типах значений и их структуре можно узнать из официальной
// документации https://www.zabbix.com/documentation/current/en/manual/api/reference/action/delete.
// В сигнатуре функции перечень id является идентификатор ранее установленных тригеров
func (api *ZabbixConnectionJsonRPC) ActionDelete(ctx context.Context, id ...string) ([]byte, error) {
	if len(id) == 0 {
		return nil, errors.New("deletion cannot be performed, the list of IDs must not be empty")
	}

	return api.sendRequest(
		ctx,
		strings.NewReader(
			fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"action.delete",
	  			"params":[%s],
	  			"id":1
			}`, strings.Join(id, ","))))
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
func (api *ZabbixConnectionJsonRPC) CustomRequest(ctx context.Context, method string, param string) ([]byte, error) {
	return api.sendRequest(
		ctx,
		strings.NewReader(fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"%s",
	  			"params": %s,
	  			"id":100
			}`, method, param)))
}

// GetFullHostList весь список хостов
func (api *ZabbixConnectionJsonRPC) GetFullHostList(ctx context.Context) ([]byte, error) {
	return api.sendRequest(
		ctx,
		strings.NewReader(`{
	  			"jsonrpc":"2.0",
	  			"method":"host.get",
	  			"params": {
					"output":"extend"
				},
	  			"id":1
			}`))
}

// GetGetHosts список хостов
func (api *ZabbixConnectionJsonRPC) GetHosts(ctx context.Context) ([]byte, error) {
	return api.sendRequest(
		ctx,
		strings.NewReader(`{
	  			"jsonrpc":"2.0",
	  			"method":"host.get",
	  			"params": {
					"output":"extend"
				},
	  			"id":1
			}`))
}

// GetHostLis список хостов для определённых групп
func (api *ZabbixConnectionJsonRPC) GetHostList(ctx context.Context, groupId ...string) ([]byte, error) {
	if len(groupId) == 0 {
		return nil, errors.New("value 'groupId' is not be empty")
	}

	return api.sendRequest(
		ctx,
		strings.NewReader(fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"host.get",
	  			"params": {
					"output":"extend",
					"selectTags": "extend",
					"selectMacros": "extend",
					"selectInventory": "extend",
					"selectInterfaces": "extend",
					"groupids": [%s]
				},
	  			"id":1
			}`, strings.Join(groupId, ","))))
}

// GetHostGroups список групп к которым принадлежит хост с заданным id
func (api *ZabbixConnectionJsonRPC) GetHostGroups(ctx context.Context, hostId string) ([]Group, error) {
	errMsg := &ResponseError{}
	res := &struct {
		Result []struct {
			HostGroups []Group `json:"hostgroups"`
			HostId     string  `json:"hostid"`
		} `json:"result"`
		JsonRPC string `json:"jsonrpc"`
		ID      int    `json:"id"`
	}{}

	//получаем информацию о хосте
	b, err := api.CustomRequest(
		ctx,
		"host.get",
		fmt.Sprintf(`{
			"hostids": "%s",
	        "selectHostGroups": "extend",
	        "evaltype": 0,
			"groups": []
			}`, hostId),
	)
	if err != nil {
		return nil, err
	}

	res, errMsg, err = supportingfunctions.ResponseUnmarchal(b, res, errMsg)
	if err != nil {
		return nil, err
	}

	//если нет ошибок, но ответ попрежнему пустой
	if len(res.Result) == 0 && errMsg.Error.Message != "" {
		return nil, errors.New(errMsg.Error.Message)
	}

	finalyResponse := []Group(nil)
	for _, v := range res.Result {
		finalyResponse = append(finalyResponse, v.HostGroups...)
	}

	return finalyResponse, nil
}

// GetFullInformationAboutHost полная информация о хосте
func (api *ZabbixConnectionJsonRPC) GetFullInformationAboutHost(ctx context.Context, hostId string) (*ResponseHostList, error) {
	res := &ResponseHostList{}
	errMsg := &ResponseError{}

	// получаем полную информацию о хосте
	b, err := api.CustomRequest(
		ctx,
		"host.get",
		fmt.Sprintf(`{
			"hostids": "%s",
	        "selectTags": "extend",
	        "selectMacros": "extend",
			"selectInventory": "extend",
	        "selectInterfaces": "extend",
	        "evaltype": 0
			}`, hostId),
	)
	if err != nil {
		return res, err
	}

	res, errMsg, err = supportingfunctions.ResponseUnmarchal(b, res, errMsg)
	if err != nil {
		return nil, err
	}

	//если нет ошибок, но ответ попрежнему пустой, проверяем на наличие ошибки в ответе
	if len(res.Result) == 0 && errMsg.Error.Message != "" {
		return nil, errors.New(errMsg.Error.Message)
	}

	if len(res.Result) == 0 {
		return res, fmt.Errorf("the host with id '%s' was not found", hostId)
	}

	return res, nil
}

// GetHostTags список тегов у определённого хоста
func (api *ZabbixConnectionJsonRPC) GetHostTags(ctx context.Context, hostId string) ([]Tag, error) {
	res, err := api.GetFullInformationAboutHost(ctx, hostId)
	if err != nil {
		return nil, err
	}

	if len(res.Result) == 0 {
		return nil, fmt.Errorf("the host with id '%s' was not found", hostId)
	}

	finalyResponse := []Tag(nil)
	for _, v := range res.Result {
		finalyResponse = append(finalyResponse, v.Tags...)
	}

	return finalyResponse, nil
}

// GetHostMacros список макросов определённого хоста
func (api *ZabbixConnectionJsonRPC) GetHostMacros(ctx context.Context, hostId string) ([]Macro, error) {
	res, err := api.GetFullInformationAboutHost(ctx, hostId)
	if err != nil {
		return nil, err
	}

	if len(res.Result) == 0 {
		return nil, fmt.Errorf("the host with id '%s' was not found", hostId)
	}

	finalyResponse := []Macro(nil)
	for _, v := range res.Result {
		finalyResponse = append(finalyResponse, v.Macros...)
	}

	return finalyResponse, nil
}

// GetHostInterface список интерфейсов хоста
func (api *ZabbixConnectionJsonRPC) GetHostInterface(ctx context.Context, hostId string) ([]ResponseInterface, error) {
	res, err := api.GetFullInformationAboutHost(ctx, hostId)
	if err != nil {
		return nil, err
	}

	if len(res.Result) == 0 {
		return nil, fmt.Errorf("the host with id '%s' was not found", hostId)
	}

	finalyResponse := []ResponseInterface(nil)
	for _, v := range res.Result {
		for _, vv := range v.Interfaces {
			finalyResponse = append(finalyResponse, vv)
		}
	}

	return finalyResponse, nil
}

// GetHostInventory данные по инвенторизации хоста
func (api *ZabbixConnectionJsonRPC) GetHostInventory(ctx context.Context, hostId string) (*ResponseInventory, error) {
	errMsg := &ResponseError{}
	res := &ResponseInventory{}

	//получаем информацию о хосте
	b, err := api.CustomRequest(
		ctx,
		"host.get",
		fmt.Sprintf(`{
			"hostids": "%s",
	        "selectInventory": "extend",
	        "evaltype": 0
			}`, hostId),
	)
	if err != nil {
		return res, err
	}

	res, errMsg, err = supportingfunctions.ResponseUnmarchal(b, res, errMsg)
	if err != nil {
		return res, err
	}

	//если нет ошибок, но ответ попрежнему пустой
	if len(res.Result) == 0 && errMsg.Error.Message != "" {
		return nil, errors.New(errMsg.Error.Message)
	}

	return res, nil
}

// GetFullHostGroupList весь список групп хостов
func (api *ZabbixConnectionJsonRPC) GetFullHostGroupList(ctx context.Context) ([]byte, error) {
	return api.sendRequest(
		ctx,
		strings.NewReader(`{
	  			"jsonrpc":"2.0",
	  			"method":"hostgroup.get",
	  			"params": {
					"output":"extend"
				},
	  			"id":1
			}`))
}

// GetHostGroup получение группы хостов по фильтру
func (api *ZabbixConnectionJsonRPC) GetHostGroup(ctx context.Context, filter FilterHostGroup) ([]byte, error) {
	b, err := json.Marshal(&filter)
	if err != nil {
		return nil, err
	}

	return api.sendRequest(ctx, strings.NewReader(fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"hostgroup.get",
	  			"params": {
					"output":"extend",
					"filter": %s
				},
	  			"id":1
			}`, string(b))))
}

// CreateHostGroup создание группы хостов (обязательны права супер-администратора)
// Подробное описание параметров https://www.zabbix.com/documentation/current/en/manual/api/reference/hostgroup/create
func (api *ZabbixConnectionJsonRPC) CreateHostGroup(ctx context.Context, name string) ([]byte, error) {
	return api.sendRequest(
		ctx,
		strings.NewReader(fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"hostgroup.create",
	  			"params": {
					"name":"%s" 
				},
	  			"id":1
			}`, name)))
}

// CreateHost создание хоста (обязательны права супер-администратора)
// Подробное описание параметров https://www.zabbix.com/documentation/current/en/manual/api/reference/host/create
// Стоит обратить внимание, что при создании хоста, у создаваемого хоста
// может быть только один основной интерфейс. Два основных интерфейса быть не может.
func (api *ZabbixConnectionJsonRPC) CreateHost(ctx context.Context, opt CreateHostOptionsRequest) ([]byte, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(&opt)
	if err != nil {
		return nil, err
	}

	if opt.Tags == nil {
		opt.Tags = make([]Tag, 0)
	}

	if opt.Groups == nil {
		opt.Groups = make([]Group, 0)
	}

	if opt.Macros == nil {
		opt.Macros = make([]Macro, 0)
	}

	if opt.Templates == nil {
		opt.Templates = make([]Template, 0)
	}

	b, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}

	return api.sendRequest(
		ctx,
		strings.NewReader(fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"host.create",
	  			"params": %s,
	  			"id":1
			}`, string(b))))
}

// UpdateHostGroup обновление группы хостов (обязательны права супер-администратора)
func (api *ZabbixConnectionJsonRPC) UpdateHostGroup(ctx context.Context, groupId, name string) ([]byte, error) {
	return api.sendRequest(
		ctx,
		strings.NewReader(fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"hostgroup.update",
	  			"params": {
					"groupid":"%s",
					"name":"%s" 
				},
	  			"id":1
			}`, groupId, name)))
}

// CreateHostParameterInventory создание инвенторизации хоста (обязательны права супер-администратора)
func (api *ZabbixConnectionJsonRPC) CreateHostParameterInventory(ctx context.Context, hostId string, req HostInventory) ([]byte, error) {
	b, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	return api.sendRequest(
		ctx,
		strings.NewReader(fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"host.update",
	  			"params": {
					"hostid":"%s",
					"inventory": %s 
				},
	  			"id":1
			}`, hostId, string(b))))

}

// UpdateHostParameterGroups обновление в хосте параметра 'группы хостов' (обязательны права супер-администратора)
func (api *ZabbixConnectionJsonRPC) UpdateHostParameterGroups(ctx context.Context, hostId string, hostGroupId string) ([]byte, error) {
	res, err := api.GetHostGroups(ctx, hostId)
	if err != nil {
		return nil, err
	}

	req := []struct {
		GroupId string `json:"groupid"`
	}{{GroupId: hostGroupId}}
	//дополняем запрос уже имеющимеся в хосте тегами
	for _, v := range res {
		req = append(req, struct {
			GroupId string `json:"groupid"`
		}{GroupId: v.GroupId})
	}

	b, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	return api.sendRequest(
		ctx,
		strings.NewReader(fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"host.update",
	  			"params": {
					"hostid":"%s",
					"groups": %s 
				},
	  			"id":1
			}`, hostId, string(b))))
}

// UpdateHostParameterTags обновление в хосте параметра 'теги' (обязательны права супер-администратора)
// Если имя найденного тега совпадает с именем принятого, для обработки, тега, то значение найденного тега
// будет переписано значением принятого тега.
func (api *ZabbixConnectionJsonRPC) UpdateHostParameterTags(ctx context.Context, hostId string, opt Tags) ([]byte, error) {
	tags, err := api.GetHostTags(ctx, hostId)
	if err != nil {
		return nil, err
	}

	for _, v := range opt.Tag {
		index := slices.IndexFunc(tags, func(tag Tag) bool {
			return tag.Tag == v.Tag
		})

		if index != -1 {
			tags[index] = v

			continue
		}

		tags = append(tags, v)
	}

	b, err := json.Marshal(tags)
	if err != nil {
		return nil, err
	}

	return api.sendRequest(
		ctx,
		strings.NewReader(fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"host.update",
	  			"params": {
					"hostid":"%s",
					"tags": %s 
				},
	  			"id":1
			}`, hostId, string(b))))
}

// UpdateHostParameterMacro обновление в хосте параметра 'макросы' (обязательны права супер-администратора)
func (api *ZabbixConnectionJsonRPC) UpdateHostParameterMacro(ctx context.Context, hostId string, opt Macros) ([]byte, error) {
	res, err := api.GetHostMacros(ctx, hostId)
	if err != nil {
		return nil, err
	}

	opt.Macro = append(opt.Macro, res...)

	for k, macro := range opt.Macro {
		if macro.Type == "" {
			opt.Macro[k].Type = "0"
		}
	}

	b, err := json.Marshal(&opt.Macro)
	if err != nil {
		return nil, err
	}

	return api.sendRequest(
		ctx,
		strings.NewReader(fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"host.update",
	  			"params": {
					"hostid":"%s",
					"macros": %s 
				},
	  			"id":1
			}`, hostId, string(b))))
}

// UpdateHostParameterInterfaces обновление в хосте параметра 'интерфейсы' (обязательны права супер-администратора)
// Стоит обратить внимание, что при обновлении интерфейса, параметр 'main' может быть равен '1', что означает, что
// интерфейс будет основным, только если нет другого основного интерфейса. Два основных интерфейса быть не может.
func (api *ZabbixConnectionJsonRPC) UpdateHostParameterInterfaces(ctx context.Context, hostId string, opt InterfacesRequest) ([]byte, error) {
	if len(opt.Interface) == 0 {
		return nil, fmt.Errorf("parameter 'opt' cannot be empty")
	}

	res, err := api.GetHostInterface(ctx, hostId)
	if err != nil {
		return nil, err
	}

	//дополняем запрос уже имеющимеся в хосте тегами
	for _, v := range res {
		m, err := strconv.ParseInt(v.Main, 10, 32)
		if err != nil {
			m = 0
		}

		t, err := strconv.ParseInt(v.Type, 10, 32)
		if err != nil {
			t = 1
		}

		uip, err := strconv.ParseInt(v.Useip, 10, 32)
		if err != nil {
			uip = 1
		}

		opt.Interface = append(opt.Interface, InterfaceOptionsRequest{
			Details: v.Details,
			IP:      v.IP,
			DNS:     v.DNS,
			Port:    v.Port,
			HostId:  v.HostId,
			Main:    int(m),
			Type:    int(t),
			Useip:   int(uip),
		})

	}

	b, err := json.Marshal(&opt.Interface)
	if err != nil {
		return nil, err
	}

	return api.sendRequest(
		ctx,
		strings.NewReader(fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"host.update",
	  			"params": {
					"hostid":"%s",
					"interfaces": %s 
				},
	  			"id":1
			}`, hostId, string(b))))
}

// DeleteHostGroup удаление группы хостов (обязательны права супер-администратора)
func (api *ZabbixConnectionJsonRPC) DeleteHostGroup(ctx context.Context, groupId ...string) ([]byte, error) {
	return api.sendRequest(
		ctx,
		strings.NewReader(fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"hostgroup.delete",
	  			"params": [%s],
				"id":1
				}`, strings.Join(groupId, ","))))
}

// DeleteHost удаление хостов (обязательны права супер-администратора)
func (api *ZabbixConnectionJsonRPC) DeleteHost(ctx context.Context, hostId ...string) ([]byte, error) {
	return api.sendRequest(
		ctx,
		strings.NewReader(fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"host.delete",
	  			"params": [%s],
				"id":1
				}`, strings.Join(hostId, ","))))
}
