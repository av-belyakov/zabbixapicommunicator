package connectionjsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func responseDecode(data []byte) (*ResponseMessage, error) {
	res := &ResponseMessage{}
	err := json.Unmarshal(data, res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// sendRequest обрабатывает запрос
func (api *ZabbixConnectionJsonRPC) sendRequest(ctx context.Context, r *strings.Reader) ([]byte, *ResponseMessage, error) {
	res, err := api.postRequest(ctx, r)
	if err != nil {
		// при возникновении ошибки пытаемся авторизоватся повторно,
		// так как при устаревании авторизационного хеша возможно появления
		// ошибки с сообщением:
		// 'Invalid params. Session terminated, re-login, please.'
		err = api.AuthorizationStart(ctx)
		if err != nil {
			return nil, nil, err
		}

		res, err = api.postRequest(ctx, r)
		if err != nil {
			return nil, nil, err
		}

		data, err := responseDecode(res)
		if err != nil {
			return res, data, err
		}
	}

	data, err := responseDecode(res)
	if err != nil {
		return res, data, err
	}

	return res, data, err
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
