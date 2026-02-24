package connectionjsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
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

	res, err := api.postRequest(ctx, data)
	if err != nil {
		return err
	}

	result := ZabbixAuthorizationData{}
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

		return fmt.Errorf("error authorization (%s %s)", shortMsg, fullMsg)
	}

	api.authorizationHash = result.Result

	return nil
}

// GetAuthorizationData хеш авторизации
func (api *ZabbixConnectionJsonRPC) GetAuthorizationData() string {
	return api.authorizationHash
}

// Logout завершение сеанса авторизации
func (api *ZabbixConnectionJsonRPC) Logout(ctx context.Context) (bool, error) {
	data := strings.NewReader(`{
	  "jsonrpc":"2.0",
	  "method":"user.logout",
	  "params": {},
	  "id":1
	}`)

	res, err := api.postRequest(ctx, data)
	if err != nil {
		return false, err
	}

	fmt.Println("Response:", string(res))

	result := struct {
		Error struct {
			Message string `json:"message"`
			Data    string `json:"data"`
		} `json:"error"`
		JsonRPC string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  bool   `json:"result"`
	}{}
	if err := json.Unmarshal(res, &result); err != nil {
		return false, err
	}

	if len(result.Error.Message) > 0 {
		return result.Result, fmt.Errorf("error %s %s", result.Error.Message, result.Error.Data)
	}

	return result.Result, nil
}
