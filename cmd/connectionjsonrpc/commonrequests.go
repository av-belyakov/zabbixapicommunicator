package connectionjsonrpc

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// sendRequest обрабатывает запрос
func (api *ZabbixConnectionJsonRPC) sendRequest(ctx context.Context, r *strings.Reader) ([]byte, error) {
	res, err := api.postRequest(ctx, r)
	if err != nil {
		// при возникновении ошибки пытаемся авторизоватся повторно,
		// так как при устаревании авторизационного хеша возможно появления
		// ошибки с сообщением:
		// 'Invalid params. Session terminated, re-login, please.'
		err = api.AuthorizationStart(ctx)
		if err != nil {
			return nil, err
		}

		return api.postRequest(ctx, r)
	}

	return res, err
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
