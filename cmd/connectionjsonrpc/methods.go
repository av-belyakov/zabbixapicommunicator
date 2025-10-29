package connectionjsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// AuthorizationStart авторизация клиента
// В результате авторизации должен быть получен хеш авторизации который будет
// добавлятся в каждый последующий запрос.
func (api *ZabbixConnectionJsonRPC) AuthorizationStart(ctx context.Context) error {
	data := strings.NewReader(fmt.Sprintf(`{
	  "jsonrpc":"2.0",
	  "method":"user.login",
	  "params": {
	    "username":"%s",
		"password":"%s"
	  },
	  "id":1
	}`, api.login, api.passwd))

	result := ZabbixAuthorizationData{}
	res, err := api.postRequest(ctx, data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(res, &result); err != nil {
		return err
	}

	if len(result.Error) > 0 {
		var shortMsg, fullMsg string
		for k, v := range result.Error {
			if k == "message" {
				shortMsg = fmt.Sprint(v)
			}
			if k == "data" {
				fullMsg = fmt.Sprint(v)
			}
		}

		return fmt.Errorf("error authorization, (%s %s)", shortMsg, fullMsg)
	}

	api.authorizationHash = result.Result

	return nil
}

// GetAuthorizationData хеш авторизации
func (api *ZabbixConnectionJsonRPC) GetAuthorizationData() string {
	return api.authorizationHash
}

// ActionGet извлечение данных.
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

// ActionDelete удаление действий установленных ранее.
// Подробнее о типах значений и их структуре можно узнать из официальной
// документации https://www.zabbix.com/documentation/current/en/manual/api/reference/action/delete.
// В сигнатуре функции перечень id является идентификатор ранее установленных тригеров
func (api *ZabbixConnectionJsonRPC) ActionDelete(ctx context.Context, id ...string) ([]byte, error) {
	if len(id) == 0 {
		return nil, errors.New("deletion cannot be performed, the list of IDs must not be empty")
	}

	res, err := api.sendRequest(
		ctx,
		strings.NewReader(
			fmt.Sprintf(`{
	  			"jsonrpc":"2.0",
	  			"method":"action.delete",
	  			"params":[%s],
	  			"id":1
			}`, strings.Join(id, ","))))

	return res, err
}

// sendRequest обрабатывает запрос
func (api *ZabbixConnectionJsonRPC) sendRequest(ctx context.Context, r *strings.Reader) ([]byte, error) {
	response, err := api.postRequest(ctx, r)
	if err != nil {
		// при возникновении ошибки пытаемся авторизоватся повторно,
		// так как при устаревании авторизационного хеша возможно появления
		// ошибки с сообщением:
		// 'Invalid params. Session terminated, re-login, please.'
		err = api.AuthorizationStart(ctx)
		if err != nil {
			return nil, err
		}

		response, err = api.postRequest(ctx, r)
		if err != nil {
			return nil, err
		}
	}

	return response, err
}

// postRequest выполняет POST запрос
func (api *ZabbixConnectionJsonRPC) postRequest(ctx context.Context, data *strings.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", api.url, data)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api.authorizationHash))
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
