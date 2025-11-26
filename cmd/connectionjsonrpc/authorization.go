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

		return fmt.Errorf("error authorization (%s %s)", shortMsg, fullMsg)
	}

	api.authorizationHash = result.Result

	return nil
}

// GetAuthorizationData хеш авторизации
func (api *ZabbixConnectionJsonRPC) GetAuthorizationData() string {
	return api.authorizationHash
}
