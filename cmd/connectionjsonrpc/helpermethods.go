package connectionjsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// GetMethodRequest запрос методом GET
func (rs *RequiestSensorInfo) GetMethodRequest(ctx context.Context, params string) (string, error) {
	request := fmt.Sprintf(`{
      "jsonrpc":"2.0",
	  "method":"item.get",
	  "params":%s,
	  "id":1
	}`, params)

	return rs.SendRequest(ctx, request)
}

// sendRequest передача запроса к API
func (sid *RequiestSensorInfo) SendRequest(ctx context.Context, str string) (string, error) {
	res, err := sid.zabbixConnection.PostRequest(ctx, strings.NewReader(str))
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
