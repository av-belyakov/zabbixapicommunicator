package connectionjsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// AuthorizationStart попытка получить хеш авторизации который необходим для
// дальнейшей работы с API
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

/*
* Здесь надо подумать

// Request запрос к API

	func (api *RequiestSensorInfo) Request(ctx context.Context, params string) (string, error) {
		request := fmt.Sprintf(`{
	      "jsonrpc":"2.0",
		  "method":"item.get",
		  "params":%s,
		  "id":1
		}`, params)

		return api.sendRequest(ctx, request)
	}

// sendRequest передача запроса к API

	func (api *RequiestSensorInfo) sendRequest(ctx context.Context, str string) (string, error) {
		res, err := api.zabbixConnection.postRequest(ctx, strings.NewReader(str))
		if err != nil {
			return "", err
		}

		var resData ResponseData
		err = json.Unmarshal(res, &resData)
		if err != nil {
			return "", err
		}

		if len(resData.Error) > 0 {
			var msg, data string

			for k, v := range resData.Error {
				if k == "message" {
					msg = fmt.Sprint(v)
				}

				if k == "data" {
					data = fmt.Sprint(v)
				}
			}

			return "", fmt.Errorf("%s. %s", msg, data)
		}

		for _, v := range resData.Result {
			for key, value := range v {
				if key == "lastvalue" {
					return fmt.Sprint(value), nil
				}
			}
		}

		return "", nil
	}
*/
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
